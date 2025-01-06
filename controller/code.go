package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已经存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或者密码错误",
	CodeServerBusy:      "服务器繁忙",

	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效的token",
}

// Msg 返回错误码对应的错误信息
// @Summary 返回错误码对应的错误信息
// @Description 错误码可以直接调用该方法，返回该方法对应的错误信息
// @Tags 错误码相关接口
// @Security ApiKeyAuth
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
