
package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	validUser    = User{Name: "kotokz", Mail: "kotokz@me.im", Pass: "123"}
	mySigningKey = "YOCALCODE"
)

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// Run this on all requests
		// Should be moved to a proper middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Next()
	})
	router.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	})

	router.POST("/user/token", func(c *gin.Context) {
		var login Login
		val := c.Bind(&login)
		fmt.Println(login.Username + " " + login.Password )
		if val != nil {
			c.JSON(401, gin.H{"code": 401, "msg": "Both name & password are required"})
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
				c.JSON(200, gin.H{"code": 500, "msg": "Server error!"})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString})
		} else {
			c.JSON(400, gin.H{"code": 400, "msg": "Error username or password!"})
		}
	})

	router.POST("/user/balance", func(c *gin.Context) {
		token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(mySigningKey))
			return b, nil
		})
		fmt.Println(err)
		if err != nil {
			c.JSON(200, gin.H{"code": 403, "msg": err.Error()})
		} else {
			if token.Valid {
				token.Claims["balance"] = 49
				tokenString, err := token.SignedString([]byte(mySigningKey))
				if err != nil {
					c.JSON(200, gin.H{"code": 500, "msg": "Server error!"})
					return
				}

				c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString})
			} else {
				c.JSON(200, gin.H{"code": 401, "msg": "Sorry, you are not validate"})
			}
		}
	})

	router.Run(":3001")
}