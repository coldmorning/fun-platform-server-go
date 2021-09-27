package main
import (
	"log"
	"github.com/go-redis/redis/v7"
	"github.com/gin-gonic/gin"
	"fun-platform-server/config"
	"fun-platform-server/middleware"
	"fun-platform-server/auth/delivery/http"
)

var client *redis.Client
var router = gin.Default()

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


func main(){

	router := gin.New()


	router.Use(middleware.Logger())

	v1_router := router.Group("/api/v1")
    {
		v1_router.POST("login",middleware.AuthRequired,http.Login)
		v1_router.POST("logout",middleware.AuthRequired,http.Logout)
		v1_router.GET("test",middleware.AuthRequired,http.Test)	
	}
	
	
	
	router.Run(":8083")
	
}