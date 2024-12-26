package models

// 定义请求的参数结构体

// ParamSignUp 注册参数结构体；可以在原有基础上新增tag
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}
