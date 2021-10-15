package authservice

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"

	authpostresql "github.com/coldmorning/fun-platform/auth/dao/postgresql"
	"github.com/coldmorning/fun-platform/model"
)

var (
	Refresh_SECRET = []byte("qweasdghtyhq")
	ACESS_SECRET   = []byte("playeresssd")

	ErrInvalidToken     = errors.New("token is invalid")
	ErrExpiredToken     = errors.New("token has expired")
	ErrSignMethodToken  = errors.New("unexpected signing method")
	ErrorMalformedToken = errors.New("token is malformed (Not correct format)")
	ErrorOtherToken     = errors.New("Couldn't handle this token")

	ErrorDelASToken           = errors.New("Couldn't remove acess token from redisk")
	ErrorDelRFToken           = errors.New("Couldn't remove refresh token from redisk")
	ErrorGetASToken           = errors.New("Couldn't get acess token from redisk")
	ErrorGetRFToken           = errors.New("Couldn't get refresh token from redisk")
	ErrorExtractTokenMetadata = errors.New("Couldn't ExtractTokenMetadata ")
)

func FindUser(u model.User) (*model.User, error) {

	user, err := authpostresql.FetchUser(u)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("(user can not find)")
	}

	return user, nil
}

func ExtractTokenMetadata(token *jwt.Token) (*model.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, ErrorGetASToken
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, ErrorGetRFToken
		}

		return &model.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, ErrorExtractTokenMetadata
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
	if err != nil {
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

}

func CreateAesccToken(userId string, td *model.TokenDetails) (*model.TokenDetails, error) {
	var err error
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = time.Now().Add(time.Second * 20).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	td.AccessToken, err = at.SignedString(ACESS_SECRET)

	if err != nil {
		return nil, err
	}
	return td, err
}

func CreateRefreshToken(userId string, td *model.TokenDetails) (*model.TokenDetails, error) {
	var err error
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + userId
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

func CreateToken(userId string) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}
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

func RemoveToken(auth *model.AccessDetails, client *redis.Client) error {
	refreshUuid := fmt.Sprintf("%s++%s", auth.AccessUuid, auth.UserId)

	access, err := client.Del(auth.AccessUuid).Result()

	if err != nil {
		return errors.New(ErrorDelASToken.Error() + ":" + err.Error())
	}
	refresh, err := client.Del(refreshUuid).Result()
	if err != nil {
		return errors.New(ErrorDelASToken.Error() + ":" + err.Error())
	}

	if access != 1 || refresh != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func CreateAuth(userId string, td *model.TokenDetails, client *redis.Client) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAcess := client.Set(td.AccessUuid, userId, at.Sub(now)).Err()
	if errAcess != nil {
		return errAcess
	}

	errRefresh := client.Set(td.RefreshUuid, userId, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func RemoveAuth(auth *model.AccessDetails, client *redis.Client) error {

	err := RemoveToken(auth, client)
	if err != nil {
		return err
	}
	return nil
}
