package http

import (
	"os"
	"log"
	"net/http"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"fun-platform-server/domain"
	"fun-platform-server/auth/usecase"

)
var client *redis.Client



func init(){
	path,err := os.Getwd()
	
	if err != nil {
		log.Println("read path error")
	}
	parent_path := path[:len(path)-4];
	
    config := viper.New()
	config.AddConfigPath(parent_path) 
	config.SetConfigType("yaml")
    config.SetConfigName ("dev-env")
	if err := config.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Config file not found; ignore error if desired
            log.Println("no such config file")
        } else {
            // Config file was found but another error was produced
            log.Println("read config error")
        }
        log.Fatal (err) // failed to read configuration file. Fatal error
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