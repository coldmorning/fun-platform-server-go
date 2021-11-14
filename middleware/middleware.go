package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/coldmorning/fun-platform/config"
	authservice "github.com/coldmorning/fun-platform/auth/service"
)

var (
	Access_secret  = []byte("")
	Refresh_secret = []byte("")
)

func init() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	Access_secret = []byte(config.GetString("Token.ACCESS_TOKEN.SECRET"))
	Refresh_secret = []byte(config.GetString("Token.REFRESH_TOKEN.SECRET"))
}

func Log(ctx *gin.Context) {
	log.Println("exec middleware1")
	ctx.Next()
	log.Println("after exec middleware1")
}

func Auth(ctx *gin.Context) {
	_, err := authservice.VerifyAccessToken(Access_secret, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers for this request are not called.
		ctx.Abort()
	}else{
		//Update JWT
		ctx.Next()
	}
}