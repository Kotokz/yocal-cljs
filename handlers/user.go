package handlers


import (

	"github.com/lunny/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
)

type User struct {
	Name string `json: "name"`
	Mail string `json: "mail"`
	Pass string `json: "pass"`
}

type Login struct {
	Username string `json: "username" binding: "required"`
	Password string `json: "password" binding: "required"`
}

var (
	validUser = User{Name: "kotokz", Mail: "kotokz@me.im", Pass: "123"}
	mySigningKey = "YOCALCODE"
)

func Token(c *gin.Context) {
	var login Login
	val := c.Bind(&login)
	log.Info(login.Username + " " + login.Password)
	if val != nil {
		c.JSON(401, gin.H{"code":http.StatusBadRequest, "msg": "Both name & password are required"})
		return
	}
	if login.Username == validUser.Name && login.Password == validUser.Pass {
		token := jwt.New(jwt.SigningMethodHS256)
		// Headers
		token.Header["alg"] = "HS256"
		token.Header["typ"] = "JWT"

		// Claims
		token.Claims["name"] = validUser.Name
		token.Claims["mail"] = validUser.Mail
		token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		tokenString, err := token.SignedString([]byte(mySigningKey))
		if err != nil {
			c.JSON(200, gin.H{"code": http.StatusInternalServerError, "msg": "Server error!"})
			return
		}
		c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK", "jwt": tokenString})
	} else {
		c.JSON(400, gin.H{"code": http.StatusUnauthorized, "msg": "Error username or password!"})
	}
}

func Balance(c *gin.Context) {
	token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
		b := ([]byte(mySigningKey))
		return b, nil
	})
	if err != nil {
		log.Errorf("MSG: %v", err)
		c.JSON(200, gin.H{"code": http.StatusForbidden, "msg": err.Error()})
		return
	}

	if token.Valid {
		token.Claims["balance"] = 49
		tokenString, err := token.SignedString([]byte(mySigningKey))
		if err != nil {
			c.JSON(200, gin.H{"code": http.StatusInternalServerError, "msg": "Server error!"})
			return
		}

		c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK", "jwt": tokenString})
	} else {
		c.JSON(200, gin.H{"code": http.StatusUnauthorized, "msg": "Sorry, you are not authorized"})
	}

}
