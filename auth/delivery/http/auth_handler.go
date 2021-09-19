package http

import (
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"fun-platform-server/domain"
	"fun-platform-server/auth/usecase"

)
var client *redis.Client



func init(){
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6380"
	}
	client = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "password123", // no password set
		DB:       0,  // use default DB
	})
	_,err := client.Ping().Result()
	if err != nil{
		panic(err)	
	}
}

func Login(c *gin.Context){

	var body domain.User
	var err error
	if err = c.ShouldBindJSON(&body);err != nil{
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := usecase.FindUser(body)
	if err != nil{
		c.JSON(http.StatusUnauthorized,"Please provide valid login details "+err.Error())
		return
	}

	token,err := usecase.CreateToken(user.ID)
	if err !=nil{
		c.JSON(http.StatusUnprocessableEntity,err.Error())
		return
	}
	err = usecase.CreateAuth(user.ID,token,client)
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity,err.Error())
	}
	tokens := map[string]string{
		"access_token" : token.AccessToken,
		"refresh_token" : token.RefreshToken,
	}
	
	c.JSON(http.StatusOK, tokens)
}

func Test(c *gin.Context){
		token, err :=usecase.VerifyToken(c.Request)
		if err !=nil{
			c.JSON(http.StatusUnprocessableEntity,err.Error())
			return
		}
		c.JSON(http.StatusOK, token)		
	
}

func Logout(C *gin.Context){
	//var u domain.User
	//var err error
	
}