package main
import (
	"log"
	"os"

	authttp "github.com/coldmorning/fun-platform/auth/controller/http"
	boardhttp "github.com/coldmorning/fun-platform/board/controller/http"
	"github.com/coldmorning/fun-platform/config"
	"github.com/coldmorning/fun-platform/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client
var router = gin.Default()

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

func main() {
	path, err := os.Getwd()

	if err != nil {
		log.Println("read path error")
	}
	swagger_path := path + "\\swagger-ui\\dist"
	router := gin.New()
	router.Static("/swaggerui", swagger_path)

	v1_router := router.Group("/api/v1")
	{
		v1_router.POST("home", middleware.Log, authttp.Home)
		v1_router.POST("login", middleware.Log, authttp.Login)
		v1_router.POST("logout", middleware.Log, authttp.Logout)
		v1_router.POST("refresh", middleware.Log, authttp.Refresh)
		v1_router.GET("test", middleware.Log, authttp.Test)
		
		boardRouter := v1_router.Group("board")
		boardRouter.Use(middleware.Auth)
		{
			boardRouter.GET("/", middleware.Auth, boardhttp.List)
			boardRouter.POST("/:id", middleware.Auth, boardhttp.Create)
			boardRouter.DELETE("/:id", middleware.Auth, boardhttp.Delete)
			boardRouter.PUT("/:id", middleware.Auth, boardhttp.Update)
			boardRouter.PATCH("/:id/state", middleware.Auth, boardhttp.Update)

		}
	}

	router.Run(":8083")

}
