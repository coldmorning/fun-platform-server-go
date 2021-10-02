package authservice

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
	
	"github.com/coldmorning/fun-platform/auth/repository/postgresql"
	"github.com/coldmorning/fun-platform/domain"
)

var (
	Refresh_SECRET = []byte("qweasdghtyhq")
	ACESS_SECRET   = []byte("playeresssd")

	ErrInvalidToken     = errors.New("token is invalid")
	ErrExpiredToken     = errors.New("token has expired")
	ErrSignMethodToken  = errors.New("unexpected signing method")
	ErrorMalformedToken = errors.New("token is malformed (Not correct format)")
	ErrorOtherToken     = errors.New("Couldn't handle this token")
)

func FindUser(u domain.User) (*domain.User, error) {

	user, err := authpostresql.FetchUser(u)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("(user can not find)")
	}

	return user, nil
}
func VerifyToken(r *http.Request) (*jwt.Token, error) {

	var tokenStr string = ""
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		tokenStr = strArr[1]
	} else {
		return nil, ErrorMalformedToken
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrSignMethodToken
		}

		return ACESS_SECRET, nil
	})
	if err == nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, ErrorMalformedToken

		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, ErrExpiredToken
		} else {
			return nil, ErrorOtherToken
		}
	} else {
		return nil, ErrorOtherToken
	}

	return token, nil
}

func CreateAesccToken(userId string, td *domain.TokenDetails) (*domain.TokenDetails, error) {
	var err error
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AcessUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["acess_uuid"] = td.AcessUuid
	atClaims["exp"] = time.Now().Add(time.Second * 20).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	td.AccessToken, err = at.SignedString(ACESS_SECRET)

	if err != nil {
		return nil, err
	}
	return td, err
}

func CreateRefreshToken(userId string, td *domain.TokenDetails) (*domain.TokenDetails, error) {
	var err error
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	td.RefreshToken, err = rt.SignedString(Refresh_SECRET)

	if err != nil {
		return nil, err
	}
	return td, err
}

func CreateToken(userId string) (*domain.TokenDetails, error) {
	td := &domain.TokenDetails{}
	var err error

	td, err = CreateAesccToken(userId, td)
	if err != nil {
		return nil, err
	}
	td, err = CreateRefreshToken(userId, td)

	if err != nil {
		return nil, err
	}

	return td, nil
}

func CreateAuth(userId string, td *domain.TokenDetails, client *redis.Client) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAcess := client.Set(td.AcessUuid, userId, at.Sub(now)).Err()
	if errAcess != nil {
		return errAcess
	}

	errRefresh := client.Set(td.RefreshUuid, userId, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
