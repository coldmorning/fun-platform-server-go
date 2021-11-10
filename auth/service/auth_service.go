package authservice

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	authpostresql "github.com/coldmorning/fun-platform/auth/dao/postgresql"
	"github.com/coldmorning/fun-platform/config"
	"github.com/coldmorning/fun-platform/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
)

var (
	ErrInvalidToken     = errors.New("token is invalid")
	ErrExpiredToken     = errors.New("token has expired")
	ErrSignMethodToken  = errors.New("unexpected signing method")
	ErrorMalformedToken = errors.New("token is malformed (Not correct format)")
	ErrorOtherToken     = errors.New("Couldn't handle this token")

	ErrorDelASTokenRedis      = errors.New("Couldn't remove access token from redisk")
	ErrorDelRFTokenRedis      = errors.New("Couldn't remove refresh token from redisk")
	ErrorGetASTokenRedis      = errors.New("Couldn't get access token from redisk")
	ErrorGetRFTokenRedis      = errors.New("Couldn't get refresh token from redisk")
	ErrorGetASToken           = errors.New("Couldn't extract access token")
	ErrorGetRFToken           = errors.New("Couldn't extract refresh token")
	ErrorGetASUuid            = errors.New("Couldn't extract access uuid")
	ErrorGetRFUuid            = errors.New("Couldn't extract refresh uuid")
	ErrorGetUserId            = errors.New("Couldn't extract  userId")
	ErrorExtractTokenMetadata = errors.New("Couldn't ExtractTokenMetadata ")
)
var (
	Access_secret  = []byte("")
	Refresh_secret = []byte("")
	Access_time    = -1
	Refresh_time   = -1
)

func init() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	Access_secret = []byte(config.GetString("Token.ACCESS_TOKEN.SECRET"))
	Refresh_secret = []byte(config.GetString("Token.REFRESH_TOKEN.SECRET"))
	Access_time = config.GetInt("Token.ACCESS_TOKEN.TIME")
	Refresh_time = config.GetInt("Token.REFRESH_TOKEN.TIME")
}
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

func ExtractAccessTokenMetadata(token *jwt.Token) (*model.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, ErrorGetASToken
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, ErrorGetUserId
		}

		return &model.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, ErrorExtractTokenMetadata
}

func ExtractRefreshTokenMetadata(token *jwt.Token) (*model.RefreshDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, ErrorGetRFUuid
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, ErrorGetUserId
		}

		return &model.RefreshDetails{
			RefreshUuid: refreshUuid,
			UserId:      userId,
		}, nil
	}
	return nil, ErrorExtractTokenMetadata
}

func VerifyAccessToken(screctKey []byte, r *http.Request) (*jwt.Token, error) {
	var tokenStr string = ""
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		tokenStr = strArr[1]
	} else {
		return nil, ErrorMalformedToken
	}
	return VerifyToken(screctKey, tokenStr)
}

func VerifyToken(screctKey []byte, tokenStr string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrSignMethodToken
		}

		return screctKey, nil
	})


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


func CreateAccessToken(userId string, td *model.TokenDetails) (*model.TokenDetails, error) {
	var err error
	td.AtExpires = time.Now().Add(time.Duration(Access_time) * time.Microsecond).Unix()
	td.AccessUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	td.AccessToken, err = at.SignedString(Access_secret)

	if err != nil {
		return nil, err
	}
	return td, err
}

func CreateRefreshToken(userId string, td *model.TokenDetails) (*model.TokenDetails, error) {
	var err error
	td.RtExpires = time.Now().Add(time.Duration(Refresh_time) * time.Microsecond).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + userId
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	td.RefreshToken, err = rt.SignedString(Refresh_secret)

	if err != nil {
		return nil, err
	}
	return td, err
}


func CreateToken(userId string) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}
	var err error

	td, err = CreateAccessToken(userId, td)
	if err != nil {
		return nil, err
	}
	td, err = CreateRefreshToken(userId, td)

	if err != nil {
		return nil, err
	}

	return td, nil
}


func RemoveTokens(auth *model.AccessDetails, client *redis.Client) error {
	refreshUuid := fmt.Sprintf("%s++%s", auth.AccessUuid, auth.UserId)

	access, err := client.Del(auth.AccessUuid).Result()

	if err != nil {
		return errors.New(ErrorDelASTokenRedis.Error() + ":" + err.Error())
	}
	refresh, err := client.Del(refreshUuid).Result()
	if err != nil {
		return errors.New(ErrorDelASTokenRedis.Error() + ":" + err.Error())
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

func RemoveAuth(givenUuid string, client *redis.Client) (int64, error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
