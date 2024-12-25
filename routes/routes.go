package routes

import (
	"Dandelion/controller"
	"Dandelion/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)

	return r
}
