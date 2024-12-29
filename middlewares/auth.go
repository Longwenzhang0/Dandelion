package middlewares

import (
	"Dandelion/controller"
	"Dandelion/pkg/jwt"
	"errors"
	"strings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头Authorization: Bearer xxx.xxx.xxx
		// 这里的具体实现方式要依据实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			zap.L().Error("JWTAuthMiddleware token error: ", zap.Error(errors.New("no authHeader")))
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			zap.L().Error("JWTAuthMiddleware token error: ", zap.Error(errors.New("invalid token format")))
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		myClaim, err := jwt.ParseToken(parts[1])
		if err != nil {
			zap.L().Error("JWTAuthMiddleware token error: ", zap.Error(err))
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, myClaim.Username)
		c.Next() // 后续的处理函数可以用过c.Get(CtxUserIDKey)来获取当前请求的用户信息
	}
}
