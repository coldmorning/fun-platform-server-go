package auth

import (
	"os"
	"time"
	"strconv"
    "github.com/dgrijalva/jwt-go"
    "github.com/twinj/uuid"
	"fun-platform-server/domain"
	"github.com/go-redis/redis/v7"
)
func CreateToken(userId uint64)(*domain.TokenDetails,error){
	td:= &domain.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute *15).Unix()
	td.AcessUuid = uuid.NewV4().String()
	
	td.RtExpires = time.Now().Add(time.Hour*24*7).Unix()
	td.RefreshUuid = uuid.NewV4().String()
	
	var err error
	os.Setenv("ACESS_SECRET","playeresssd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["acess_uuid"] = td.AcessUuid
	atClaims["exp"] = time.Now().Add(time.Second *20).Unix()
	at:=jwt.NewWithClaims(jwt.SigningMethodHS512,atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACESS_SECRET")))
	if err != nil {
		return nil,err
	}
	os.Setenv("Refresh_SECRET","qweasdghtyhq")
	rtClaims:= jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt :=jwt.NewWithClaims(jwt.SigningMethodHS512,rtClaims)
	td.RefreshToken,err = rt.SignedString([]byte(os.Getenv("Refresh_SECRET")))
	if err != nil{
		return nil,err
	}

	return td, nil
}

func CreateAuth(userId uint64, td *domain.TokenDetails, client *redis.Client) error{
	at := time.Unix(td.AtExpires,0)
	rt := time.Unix(td.RtExpires,0)
	now := time.Now()

	errAcess := client.Set(td.AcessUuid,  strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAcess != nil{
		return errAcess
	}

	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil{
		return errRefresh
	}
	return nil
}