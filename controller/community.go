package controller

import (
	"Dandelion/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关
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
