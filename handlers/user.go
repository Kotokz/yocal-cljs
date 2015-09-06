package handlers

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/kotokz/yocal-cljs/modules/forms"
)

type User struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
	Pass string `json:"pass"`
}

var (
	validUser    = User{Name: "kotokz", Mail: "kotokz@me.im", Pass: "123"}
	MySigningKey = "YOCALCODE"
)

func Token(c *gin.Context) {
	var login forms.LoginForm
	val := c.Bind(&login)

	if val != nil {
//		var loginError forms.LoginForm
		out := forms.ParseFormErrors(val)

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "input format incorrect",
			"errors":out})
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
	var reg forms.RegisterForm
	val := c.Bind(&reg)
	if val != nil {
		//		var loginError forms.LoginForm
		out := forms.ParseFormErrors(val)

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "input format incorrect",
			"errors":out})
		return
	}

	log.Info(reg)
	c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK"})
}


func Balance(c *gin.Context) {
	c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK", "balance": 49})
}
