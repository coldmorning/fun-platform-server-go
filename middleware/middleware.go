package middleware

import (
	"time"
	"log"
	"github.com/gin-gonic/gin"

)



func AuthRequired(c *gin.Context)  {
	log.Println("exec middleware1")
	c.Next()
	log.Println("after exec middleware1")
}