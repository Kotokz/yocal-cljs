package handlers

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type User struct {
	Name string `json: "name"`
	Mail string `json: "mail"`
	Pass string `json: "pass"`
}

type Login struct {
	Username string `json:"username" binding:"required,min=7,max=7"`
	Password string `json:"password" binding:"required,min=7"`
}

var (
	validUser    = User{Name: "kotokz", Mail: "kotokz@me.im", Pass: "123"}
	MySigningKey = "YOCALCODE"
)

func Token(c *gin.Context) {
	var login Login
	val := c.Bind(&login)

	if val != nil {
		log.Info(val.Error())
		c.JSON(401, gin.H{"code": http.StatusBadRequest, "msg": val.Error()})
		return
	}
	log.Info("user name = "+login.Username)
	if login.Username == validUser.Name && login.Password == validUser.Pass {
		token := jwt.New(jwt.SigningMethodHS256)
		// Headers
		token.Header["alg"] = "HS256"
		token.Header["typ"] = "JWT"

		// Claims
		token.Claims["name"] = validUser.Name
		token.Claims["mail"] = validUser.Mail
		token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		tokenString, err := token.SignedString([]byte(MySigningKey))
		if err != nil {
			c.JSON(200, gin.H{"code": http.StatusInternalServerError, "msg": "Server error!"})
			return
		}
		c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK", "jwt": tokenString})
	} else {
		c.JSON(400, gin.H{"code": http.StatusUnauthorized, "msg": "Error username or password!"})
	}
}

func Register(c *gin.Context) {

}


func Balance(c *gin.Context) {
	c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK", "balance": 49})
}
