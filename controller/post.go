package controller

import (
	"Dandelion/logic"
	"Dandelion/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 发布帖子handler
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数和参数校验；创建实例，将参数映射到实例；结构体按需加tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreatePostHandler/c.ShouldBindJSON(p) failed: ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从ctx中取到当前的用户ID；auth之后就存到了ctx中，使用getCurrentUserID获取
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("CreatePostHandler/getCurrentUserID(c) failed: ", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("CreatePostHandler/logic.CreatePost(p) failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetail 通过id获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数和参数校验，帖子id
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetailHandler/strconv.ParseInt(pidStr,10,64) failed: ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据id取出帖子数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostDetailHandler/logic.GetPostByID(pid) failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}
