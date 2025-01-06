package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	// Msg和Data可能是多种类型，所以这里用泛型接收
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 根据code返回错误信息
// @Summary 根据code返回错误信息
// @Description 根据code返回错误信息
// @Tags 错误码相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseErrorWithMsg 根据code和msg返回错误信息
// @Summary 根据code和msg返回错误信息
// @Description 根据code和msg返回错误信息
// @Tags 错误码相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {

	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseSuccess 返回成功响应
// @Summary 返回成功响应
// @Description 返回成功响应
// @Tags 错误码相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200
func ResponseSuccess(c *gin.Context, data interface{}) {

	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
