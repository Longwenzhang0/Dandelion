package routes

import (
	"Dandelion/controller"
	"Dandelion/logger"
	"Dandelion/middlewares"
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

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 登录才可以访问，判断请求头jwt token是否有效
		c.String(http.StatusOK, "pong!")
	})

	// 登录业务路由
	r.POST("/login", controller.LoginHandler)

	return r
}
