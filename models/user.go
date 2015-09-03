package models

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/lunny/log"

	"git.oschina.net/kotokz/yocal/modules/utils"
)

type User struct {
	Id        int64
	LowerName string `xorm:"UNIQUE NOT NULL"`
	LoginName string `xorm:"varchar(30) unique"`
	FullName  string `xorm:"varchar(30) unique"`
	Staffid   int    `xorm:"varchar(80) unique"`
	Email     string `xorm:"varchar(80) unique"`
	Password  string `xorm:"varchar(128)"`
	Role      int
	IsAdmin   bool      `xorm:"index"`
	IsActive  bool      `xorm:"index"`
	Rands     string    `xorm:"varchar(10)"`
	Salt      string    `xorm:"varchar(10)"`
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}

// EncodePasswd encodes password to safe format.
func (u *User) EncodePasswd() {
	newPasswd := utils.PBKDF2([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(passwd string) bool {
	newUser := &User{Password: passwd, Salt: u.Salt}
	newUser.EncodePasswd()
	return u.Password == newUser.Password
}

// GetUserSalt returns a ramdom user salt token.
func GetUserSalt() string {
	return utils.GetRandomString(10)
}

func CreateUser(user *User) error {

	isExist, err := IsUserExistByLoginName(user.LoginName, 0)
	if err != nil {
		return err
	} else if isExist {
		return ErrNameAlreadyExist
	}

	isExist, err = IsUserExistByEmail(user.Email, 0)
	if err != nil {
		return err
	} else if isExist {
		return ErrEmailAlreadyUsed
	}

	isExist, err = IsUserExistByStaffid(user.Staffid, 0)
	if err != nil {
		return err
	} else if isExist {
		return ErrStaffIdAlreadyExist
	}

	user.LowerName = strings.ToLower(user.LoginName)

	user.Rands = GetUserSalt()
	user.Salt = GetUserSalt()

	user.EncodePasswd()

	// Auto-set admin for the first user.
	c, err := Count(&User{IsAdmin: true})
	if err != nil {
		return err
	} else if c < 1 {
		user.IsActive = true
		user.IsAdmin = true
	}

	// no special activation at the moment, so all reigstered users are active

	user.IsActive = true

	return Insert(user)
}

func GetUserById(id int64) (*User, error) {
	var user User
	err := GetById(id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByLoginName(name string) (*User, error) {
	var user = User{
		LoginName: name,
	}
	err := GetByExample(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsUserExistByLoginName(name string, skipId int64) (bool, error) {
	return orm.Where("id <> ?", skipId).Get(&User{LoginName: name})
}

func IsUserExistByEmail(email string, skipId int64) (bool, error) {
	return orm.Where("id <> ?", skipId).Get(&User{Email: email})
}

func IsUserExistByStaffid(sid int, skipId int64) (bool, error) {
	return orm.Where("id <> ?", skipId).Get(&User{Staffid: sid})
}

func UserLogin(uname, passwd string) (*User, error) {
	user := new(User)

	if strings.Contains(uname, "@") {
		user = &User{Email: uname}
	} else {
		user = &User{LowerName: strings.ToLower(uname)}
	}

	if err := GetByExample(user); err != nil {
		return nil, err
	}

	log.Debugf("validating password of user %s", user.LoginName)

	if !user.ValidatePassword(passwd) {
		return nil, ErrPWDIncorrect
	}

	return user, nil
}
