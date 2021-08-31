package auth_handler

import (
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"fun-platform-server/domain"
	"fun-platform-server/auth/usercase"

)
var client *redis.Client

var user = domain.User{
	ID: 1,
	Username: "username@gmail.com",
	Password: "password",
}

func init(){
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_,err := client.Ping().Result()
	if err != nil{
		panic(err)	
	}
}

func Login(c *gin.Context){

	var u domain.User
	var err error
	if err = c.ShouldBindJSON(&u);err != nil{
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if user.Username != u.Username || user.Password !=u.Password{
		c.JSON(http.StatusUnauthorized,"Please provide valid login details")
		return
	}

	token,err := auth.CreateToken(user.ID)
	if err !=nil{
		c.JSON(http.StatusUnprocessableEntity,err.Error())
		return
	}
	err = auth.CreateAuth(user.ID,token,client)
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity,err.Error())
	}
	tokens := map[string]string{
		"access_token" : token.AccessToken,
		"refresh_token" : token.RefreshToken,
	}
	
	c.JSON(http.StatusOK, tokens)
}