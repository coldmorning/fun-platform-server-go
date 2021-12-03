package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	authservice "github.com/coldmorning/fun-platform/auth/service"
	"github.com/coldmorning/fun-platform/config"
)

var ctx = context.Background()
var (
	Access_secret  = []byte("")
	Refresh_secret = []byte("")
	client         *redis.Client
	router         = gin.Default()
)

func init() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	Access_secret = []byte(config.GetString("Token.ACCESS_TOKEN.SECRET"))
	Refresh_secret = []byte(config.GetString("Token.REFRESH_TOKEN.SECRET"))
}

func Log(ctx *gin.Context) {
	log.Println("exec middleware1")
	ctx.Next()
	log.Println("after exec middleware1")
}

func Auth(ctx *gin.Context) {
	_, err := authservice.VerifyAccessToken(Access_secret, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers for this request are not called.
		ctx.Abort()
	} else {
		//Update JWT
		ctx.Next()
	}
}

func IpFilter(ctx *gin.Context) {
	requestCount := 1000
	requestIp := ctx.ClientIP()
	requestRmainInt := 0
	requestNextTime := time.Now().Add(time.Duration(1) * time.Hour).Format("2021-01-01 12:12")

	redisExpired := 1 * time.Hour
	redisKey := fmt.Sprintf("%s-%s", "IpFilter", requestIp)

	strCount, errClient := client.Get(redisKey).Result()
	currCount, err := strconv.Atoi(strCount)

	if errClient != nil {
		//fmt.Println(currCount, "is not an integer.")
	}
	if err != nil {
		_ = client.Set(redisKey, "1", redisExpired).Err()
		ctx.Header("X-RateLimit-Reset", requestNextTime)
		requestRmainInt = 999
	} else {
		currCount++
		if currCount+1 > requestCount {
			ctx.JSON(http.StatusTooManyRequests, err.Error())
		}
		requestRmainInt = requestCount - currCount
	}

	remainStr := strconv.Itoa(requestRmainInt)
	ctx.Header("X-RateLimit-Remaining", remainStr)

}
