package main

import (
	"github.com/gin-gonic/gin"

	"github.com/kotokz/yocal-cljs/handlers"
	"github.com/kotokz/yocal-cljs/modules/middlewares"
	"github.com/kotokz/yocal-cljs/modules/setting"
)

func main() {

	r := gin.New()

	r.Use(gin.Recovery(),
		gin.LoggerWithWriter(setting.LogIO))

	r.Use(func(c *gin.Context) {
		// Run this on all requests
		// Should be moved to a proper middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Next()
	})

	gin.SetMode(setting.GinMode)

	r.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	})

	user := r.Group("/user")
	{
		user.POST("/token", handlers.Token)
		user.POST("/register", handlers.Register)
	}


	authorized := r.Group("/api", middlewares.Auth(handlers.MySigningKey))

	authorized.POST("/balance", handlers.Balance)

	r.Run(setting.AppCfg.HttpPort)
}
