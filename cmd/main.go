package main
import (
	"log"
	"github.com/gin-gonic/gin"
	"fun-platform-server/auth/delivery/http"
)


var router = gin.Default()

func main(){

	router.POST("/api/v1/login",http.Login)
	router.POST("/api/v1/logout",http.Logout)
	router.GET("/api/v1/test",http.Test)
	log.Fatal(router.Run(":8080"))
}