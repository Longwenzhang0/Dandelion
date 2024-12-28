package routes

import (
	"Dandelion/controller"
	"Dandelion/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)

	// 登录业务路由
	r.POST("/login", controller.LoginHandler)

	return r
}
