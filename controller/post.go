package controller

import (
	"Dandelion/logic"
	"Dandelion/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
// @Router /post [post]
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

// GetPostDetailHandler 通过id获取帖子详情
// @Summary 通过id获取帖子详情
// @Description 通过id获取帖子详情
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} models.ApiPostDetail
// @Router /post/:id [get]
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

// GetPostListHandler 分页获取帖子列表的处理函数（mysql）
// @Summary 分页获取帖子列表的处理函数（mysql）
// @Description 分页获取帖子列表的处理函数（mysql）
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} models.ApiPostDetail
// @Router /posts/ [get]
func GetPostListHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	page, size := getPageInfo(c)

	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("GetPostListHandler/logic.GetPostList() failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)

}

// GetPostListHandler2 升级版：获取帖子列表的处理函数，redis+mysql
// @Summary 可根据前端传入参数动态获取:time or score
// @Description 可根据前端传入参数动态获取:time or score，
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} models.ApiPostDetail
// @Router /posts2/ [get]
func GetPostListHandler2(c *gin.Context) {
	// 1. 获取参数和参数校验，例如GET posts/?page=1&size=1&order=time
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2/c.ShouldBindQuery(p) failed: ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务处理，优化两个接口为一个
	data, err := logic.GetPostListSummary(p)
	if err != nil {
		zap.L().Error("GetPostListHandler2/logic.GetPostList2(p) failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, data)
}
