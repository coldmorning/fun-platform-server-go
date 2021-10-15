package main

import (
	"log"
	"os"

	authttp "github.com/coldmorning/fun-platform/auth/controller/http"
	categoryhttp "github.com/coldmorning/fun-platform/category/controller/http"
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
		v1_router.POST("login", middleware.AuthRequired, authttp.Login)
		v1_router.POST("logout", middleware.AuthRequired, authttp.Logout)
		v1_router.POST("refresh", middleware.AuthRequired, authttp.Refresh)
		v1_router.GET("test", middleware.AuthRequired, authttp.Test)

		v1_router.GET("category", middleware.AuthRequired, categoryhttp.List)
		v1_router.POST("category/:id", middleware.AuthRequired, categoryhttp.Create)
		v1_router.DELETE("category/:id", middleware.AuthRequired, categoryhttp.Delete)
		v1_router.PUT("category/:id", middleware.AuthRequired, categoryhttp.Update)
		v1_router.PATCH("category:id/state", middleware.AuthRequired, categoryhttp.Update)

	}

	router.Run(":8083")

}
