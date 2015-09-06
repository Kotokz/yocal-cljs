package forms
import (
	"strings"
	"gopkg.in/bluesuncorp/validator.v5"
	log "github.com/Sirupsen/logrus"
)

type RegisterForm struct {
	Username string `json:"username" binding:"required,min=5,max=35,excludesall= !@#?"`
	FullName  string `json:"fullname" binding:"required,min=5,max=35"`
	Staffid   int    `json:"staffid" binding:"required,number"`
	Email     string `json:"email" binding:"required,email,max=50"`
	Password  string `json:"password" binding:"required,min=6,max=255"`
	Retype    string `json:"retype" binding:"required,min=6,max=255,eqfield=Password"`
}

type LoginForm struct {
	Username string `json:"username" binding:"required,min=5,max=35,excludesall= !@#?"`
	Password  string `json:"password" binding:"required,min=6"`
	//Remember  bool   `json:"remember"`
}

func ParseFormErrors(val interface{}) map[string]string {

	if errs, ok := val.(*validator.StructErrors); ok {
		out := make(map[string]string, len(errs.Errors))
		for n, v := range errs.Errors {
			name := strings.ToLower(n)
			out[name] = v.Tag
			switch v.Tag {
			case "min":
				out[name] = "It's too short, should be more than " + v.Param + " chars"
			case "max":
				out[name] = "It's too long, should be less than " + v.Param + " chars"
			case "excludesall":
				out[name] = "It contains invalid char: " + v.Param
			case "email":
				out[name] = "It needs to be correct email format"
			case "number":
				out[name] = "It only can be number"
			case "eqfield":
				out[name] = "Should be same as " + v.Param
			default:
				out[name] = "Invalid format: " + v.Param
			}
			log.Infof("error name %v, val %v", name, v.Tag)
		}
		return out
	} else {
		out := make(map[string]string,1)
		out["other"] = "Invalid Format"
		return out
	}
}