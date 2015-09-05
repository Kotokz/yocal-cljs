
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kotokz/yocal-cljs/modules/setting"
	"github.com/kotokz/yocal-cljs/handlers"
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

	gin.SetMode(setting.GinMode)

	router.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	})

	router.POST("/user/token", handlers.Token)

	router.POST("/user/balance", handlers.Balance)

	router.Run(setting.AppCfg.HttpPort)
}