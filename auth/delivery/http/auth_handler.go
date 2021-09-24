package http

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"fun-platform-server/config"
	"fun-platform-server/domain"
	"fun-platform-server/auth/usecase"

)
var client *redis.Client



func init(){
	config,err := config.GetConfig();
	if err != nil{
		log.Fatal (err)
	}

	client = redis.NewClient(&redis.Options{
		Addr:     config.GetString("REDIS_NODE1.ENDPOINT"),
		Password: config.GetString("REDIS_NODE1.PASSWORD"), // no password set
		DB:       0,  // use default DB
	})
	_,err = client.Ping().Result()
	if err != nil{
		log.Fatal (err)
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

	token,err := usecase.CreateToken(user.Uuid)
	if err !=nil{
		c.JSON(http.StatusUnprocessableEntity,err.Error())
		return
	}
	err = usecase.CreateAuth(user.Uuid,token,client)
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