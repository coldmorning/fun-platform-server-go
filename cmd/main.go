package main
import (
	"log"
	"github.com/gin-gonic/gin"
	"fun-platform-server/auth/delivery"
)


var router = gin.Default()

func main(){
	router.POST("/api/v1/login",auth_handler.Login)
	log.Fatal(router.Run(":8080"))
}