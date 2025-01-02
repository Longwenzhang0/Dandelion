package controller

import (
	"errors"
	"strconv"

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

func getPageInfo(c *gin.Context) (int64, int64) {
	// 获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64 // 页码
		size int64 // 每页条数
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostListHandler/strconv.ParseInt(pageStr, 10, 64) failed: ", zap.Error(err))
		page = 1 // 出错时，缺省为第1页
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostListHandler/strconv.ParseInt(sizeStr, 10, 64) failed: ", zap.Error(err))
		size = 10 // 出错时，缺省为10条
	}
	return page, size
}
