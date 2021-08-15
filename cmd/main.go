package main
import (
	"os"
	"time"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"

)

type User struct {
	ID uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
func Login(c *gin.Context){
	var user =User{
		ID: 1,
		Username: "username@gmail.com",
		Password: "password",
	}
	var u User
	if err := c.ShouldBindJSON(&u);err != nil{
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if user.Username != u.Username || user.Password !=u.Password{
		c.JSON(http.StatusUnauthorized,"Please provide valid login details")
		return
	}
	token,err:= CreateToken(user.ID)
	if err !=nil{
		c.JSON(http.StatusUnprocessableEntity,err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

func CreateToken(userId uint64)(string,error){
	var err error
	os.Setenv("ACESS_SECRET","playeresssd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Second *20).Unix()
	at:=jwt.NewWithClaims(jwt.SigningMethodHS256,atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACESS_SECRET")))
	if err != nil {
		return "",err
	}
	return token, nil
}
var (
	router = gin.Default()
)
func main(){
	router.POST("/login",Login)
	log.Fatal(router.Run(":8080"))
}