package middlewares

import (
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	Bearer = "Bearer"
)

var UnAuthError = errors.New("Sorry, you are not authorized")

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.Request.Header.Get("Authorization")
		l := len(Bearer)

		if len(auth) > l+1 && auth[:l] == Bearer {
			token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
				b := ([]byte(secret))
				return b, nil
			})

			if err == nil && token.Valid {
				c.Set("claims", token.Claims)

				if exp, ok := token.Claims["exp"].(float64); ok {
					log.WithFields(log.Fields{
						"UserName": token.Claims["name"],
						"Exp":      (int64(exp) - time.Now().Unix()) / 60,
					}).Debug("User authorized")
				} else {
					log.Errorf("Incorrect claims, %v", token.Claims)
				}
			} else {
				c.AbortWithError(http.StatusUnauthorized, UnAuthError)
			}

		} else {
			//			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "Sorry, you are not authorized"})
			c.AbortWithError(http.StatusUnauthorized, UnAuthError)
		}
	}
}
