package handlers

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/kotokz/yocal-cljs/models"
	"github.com/kotokz/yocal-cljs/modules/forms"
)

type User struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
	Pass string `json:"pass"`
}

var (
	MySigningKey = "YOCALCODE"
)

func Token(c *gin.Context) {
	var form forms.LoginForm
	val := c.Bind(&form)

	if val != nil {
		//		var loginError forms.LoginForm
		out := forms.ParseFormErrors(val)

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "input format incorrect",
			"errors": out})
		return
	}

	// verify user
	user, err := models.UserLogin(form.Username, form.Password)
	if err != nil {
		c.JSON(400, gin.H{"code": http.StatusUnauthorized, "msg": "Error username or password!"})
		return
	}

	if !user.IsActive {
		c.JSON(400, gin.H{"code": http.StatusUnauthorized, "msg": "User cannot login"})
		return
	}

	log.Debugf("User logged in : %v", user.FullName)

	token := jwt.New(jwt.SigningMethodHS256)
	// Headers
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"

	// Claims
	token.Claims["name"] = user.LoginName
	token.Claims["fname"] = user.FullName
	token.Claims["mail"] = user.Email
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	tokenString, err := token.SignedString([]byte(MySigningKey))
	if err != nil {
		c.JSON(200, gin.H{"code": http.StatusInternalServerError, "msg": "Server error!"})
		return
	}
	c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK", "jwt": tokenString})
}

func Register(c *gin.Context) {
	var form forms.RegisterForm
	val := c.Bind(&form)
	if val != nil {
		//		var loginError forms.LoginForm
		out := forms.ParseFormErrors(val)

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "input format incorrect",
			"errors": out})
		return
	}

	u := &models.User{
		LoginName: form.Username,
		FullName:  form.FullName,
		Staffid:   form.Staffid,
		Email:     form.Email,
		Password:  form.Password,
	}

	if err := models.CreateUser(u); err != nil {
		var errMsg string
		switch err {
		case models.ErrNameAlreadyExist:
			errMsg = "Name already exist"
		case models.ErrEmailAlreadyUsed:
			errMsg = "Email has been taken"
		case models.ErrStaffIdAlreadyExist:
			errMsg = "staff ID has been taken"
		default:
			errMsg = err.Error()
		}
		log.Errorf("Can't save user %v, %v", u, err)

		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": errMsg})
		return
	}

	log.Info(form)
	c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK"})
}

func Balance(c *gin.Context) {
	c.JSON(200, gin.H{"code": http.StatusOK, "msg": "OK", "balance": 49})
}
