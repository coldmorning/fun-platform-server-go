package authttp

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	authservice "github.com/coldmorning/fun-platform/auth/service"
	"github.com/coldmorning/fun-platform/config"
	"github.com/coldmorning/fun-platform/model"
)

var client *redis.Client

func init() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	client = redis.NewClient(&redis.Options{
		Addr:     config.GetString("REDIS_NODE1.ENDPOINT"),
		Password: config.GetString("REDIS_NODE1.PASSWORD"), // no password set
		DB:       0,                                        // use default DB
	})
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}

func Login(ctx *gin.Context) {

	var body model.User
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := authservice.FindUser(body)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "Please provide valid login details "+err.Error())
		return
	}

	token, err := authservice.CreateToken(user.Uuid)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = authservice.CreateAuth(user.Uuid, token, client)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	ctx.JSON(http.StatusOK, tokens)
}

func Test(ctx *gin.Context) {
	token, err := authservice.VerifyToken(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, token)

}

func Logout(ctx *gin.Context) {
	//var u model.User
	//var err error
	token, err := authservice.VerifyToken(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	accessDetails, err := authservice.ExtractTokenMetadata(token)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	err = authservice.RemoveToken(accessDetails, client)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "error to remove token")
		return
	}
	ctx.JSON(http.StatusOK, "Successfully logged out")
}
