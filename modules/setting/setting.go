package setting

import (
	"io"
	"os"

	"github.com/lunny/log"
	"gopkg.in/ini.v1"
	"github.com/gin-gonic/gin"
)

const (
	CFG_PATH        = "config/yocal.ini"
)

type App struct {
	AppName   string
	AppVer    string
	AppSubUrl string
	HttpPort  string
	Https     bool
	HttpsCert string
	HttpsKey  string
}

type Db struct {
	DriverName string
	Host       string
	Port       int
	Name       string
	User       string
	Passwd     string
	MaxIdle    int
	MaxOpen    int
	DebugLog   bool
	Ssl_mode   string
	Path       string
}

type Cookie struct {
	LogInRememberDays  int
	CookieRememberName string
	CookieUserName     string
}

var (
	Cfg       *ini.File
	AppCfg    *App
	DbCfg     *Db
	CookieCfg *Cookie
// Session settings.
	Langs      []string
	Names      []string
	Redirct    bool
	LogIO      io.Writer
	GinMode    string
)

// init imports settings from yocao.ini file
func init() {
	var err error

//	f, err := os.Create("yocal.log")
//	if err != nil {
//		log.Panic("create log file failed:", err)
//	}

//	LogIO = io.MultiWriter(f, os.Stdout)
	LogIO = io.Writer(os.Stdout)
	log.SetOutput(LogIO)

	source := []interface{}{CFG_PATH}


	Cfg, err := ini.Load(source[0])
	if err != nil {
		log.Fatalf("Failed to set configuration: %v", err)
	}

	// if run_mode not prod, then marcaron.Env will use default value (dev)
	if Cfg.Section("").Key("RUN_MODE").MustString("dev") == "prod" {
		log.SetOutputLevel(log.Linfo)
		GinMode = gin.ReleaseMode
	} else {
		log.SetOutputLevel(log.Ldebug)
		GinMode = gin.DebugMode
	}

	AppCfg = &App{
		HttpPort: ":8080",
	}

	err = Cfg.Section("app").MapTo(AppCfg)
	if err != nil {
		log.Fatalf("Failed to set app conf: %v", err)
	}

	Langs = Cfg.Section("i18n").Key("LANGS").Strings(",")
	Names = Cfg.Section("i18n").Key("NAMES").Strings(",")
	Redirct = Cfg.Section("i18n").Key("REDIRECT").MustBool()

	DbCfg = &Db{
		DriverName: "mysql",
		MaxIdle:    30,
		MaxOpen:    50,
		DebugLog:   false,
	}
	err = Cfg.Section("db").MapTo(DbCfg)
	if err != nil {
		log.Fatalf("Failed to set app conf: %v", err)
	}

	CookieCfg = &Cookie{
		LogInRememberDays:  7,
		CookieRememberName: "yocal_idea",
		CookieUserName:     "yocal_think",
	}
	err = Cfg.Section("cookie").MapTo(CookieCfg)
	if err != nil {
		log.Fatalf("Failed to set app conf: %v", err)
	}

}
