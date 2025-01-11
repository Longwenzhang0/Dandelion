package routes

import (
	"Dandelion/controller"
	"Dandelion/logger"
	"Dandelion/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "Dandelion/docs" // 千万不要忘了导入把你上一步生成的docs

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	// 限流中间件，每2s放入一个令牌
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/api/v1")

	v1.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登录业务路由
	v1.POST("/login", controller.LoginHandler)
	// 应用jwt认证中间件

	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts/", controller.GetPostListHandler)
		// 根据时间或者分数获取帖子列表
		v1.GET("/posts2/", controller.GetPostListHandler2)

		// 投票
		v1.POST("/vote", controller.PostVoteController)

	}

	return r
}
