package authttp

import (
	"log"
	"net/http"
	"github.com/go-redis/redis/v7"
	"github.com/gin-gonic/gin"
	
	"github.com/coldmorning/fun-platform/config"
	"github.com/coldmorning/fun-platform/domain"
	"github.com/coldmorning/fun-platform/auth/service"
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
	user, err := authservice.FindUser(body)
	if err != nil{
		c.JSON(http.StatusUnauthorized,"Please provide valid login details "+err.Error())
		return
	}

	token,err := authservice.CreateToken(user.Uuid)
	if err !=nil{
		c.JSON(http.StatusUnprocessableEntity,err.Error())
		return
	}
	err = authservice.CreateAuth(user.Uuid,token,client)
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
		token, err :=authservice.VerifyToken(c.Request)
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