package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// 获取当前登录的用户ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	// 从ctx中获取userid,错误处理
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		zap.L().Error("getCurrentUser/c.Get(middlewares.CtxUserIDKey) failed: ", zap.Error(err))
		return
	}
	// 转换为int64
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		zap.L().Error("getCurrentUser/uid.(int64): ", zap.Error(err))
		return
	}
	return
}
