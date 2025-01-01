package controller

import (
	"Dandelion/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关

// CommunityHandler 查询所有的社区id和社区name
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区community_id,community_name，以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 查询所有的社区分类的详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id,restful传过来的
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		zap.L().Error("CommunityDetailHandler/strconv.ParseInt(communityID, 10, 64) failed: ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 业务处理，根据id查询数据
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
