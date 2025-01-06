package controller

import (
	"Dandelion/logic"
	"Dandelion/models"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// PostVoteController 给指定帖子投票接口
// @Summary 给指定帖子投票接口，获取当前登录者id，
// @Description 分页获取帖子列表的处理函数（mysql），给指定帖子投票
// @Tags 投票相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
// @Router /vote [post]
func PostVoteController(c *gin.Context) {
	// 获取参数和参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("PostVoteController/c.ShouldBindJSON(p) failed", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 类型断言，如果不是，直接返回其他报错类型；如果是，则翻译
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans))) // 翻译错误类型，去掉错误提示中的结构体标识
		return
	}
	// 获取当前登录者id，即投票者
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("PostVoteController/getCurrentUserID(c) failed: ", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 业务处理
	if err = logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("PostVoteController/logic.VoteForPost(userID, p) failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, nil)

}
