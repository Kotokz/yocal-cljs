package models

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kotokz/yocal-cljs/modules/setting"
	"os"
	"path"
	"github.com/gin-gonic/gin"
)

var (
	orm                    *xorm.Engine
	ErrNotExist            = errors.New("Not Exist")
	ErrEmailAlreadyUsed    = errors.New("Email already used")
	ErrStaffIdAlreadyExist = errors.New("Staff id already exist")
	ErrPWDIncorrect        = errors.New("Password incorrect")
	ErrNameAlreadyExist    = errors.New("Name already exist")
)

func init() {
	var err error
	cnnstr := ""

	switch setting.DbCfg.DriverName {
	case "mysql":
		if setting.DbCfg.Host[0] == '/' { // looks like a unix socket
			cnnstr = fmt.Sprintf("%s:%s@unix(%s:%d)/%s?charset=utf8&parseTime=true",
				setting.DbCfg.User, setting.DbCfg.Passwd, setting.DbCfg.Host, setting.DbCfg.Port, setting.DbCfg.Name)
		} else {
			cnnstr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
				setting.DbCfg.User, setting.DbCfg.Passwd, setting.DbCfg.Host, setting.DbCfg.Port, setting.DbCfg.Name)
		}
	case "postgres":
		cnnstr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			setting.DbCfg.User, setting.DbCfg.Passwd, setting.DbCfg.Host, setting.DbCfg.Port, setting.DbCfg.Name, setting.DbCfg.Ssl_mode)
	case "sqlite3":
		os.MkdirAll(path.Dir(setting.DbCfg.Path), os.ModePerm)
		cnnstr = "file:" + setting.DbCfg.Path + "?cache=shared&mode=rwc"
	default:
		panic(fmt.Errorf("Unknown database type: %s", setting.DbCfg.DriverName))
	}

	orm, err = xorm.NewEngine(setting.DbCfg.DriverName, cnnstr)

	if err != nil {
		panic(err)
	}

	orm.SetMaxIdleConns(setting.DbCfg.MaxIdle)
	orm.SetMaxOpenConns(setting.DbCfg.MaxOpen)

	if setting.GinMode == gin.DebugMode{
		orm.ShowSQL = true
	}

	if setting.DbCfg.DebugLog {
		orm.ShowDebug = true
		orm.ShowWarn = true
	}

	err = orm.Sync2(new(User), new(Team), new(TeamUser))
	if err != nil {
		panic(err)
	}
}
