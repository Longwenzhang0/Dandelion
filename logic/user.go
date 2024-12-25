package logic

import (
	"Dandelion/dao/mysql"
	"Dandelion/models"
	"Dandelion/pkg/snowflake"
)

// LOGIC层 用于存放业务逻辑的代码

// SignUp 注册用户
func SignUp(p *models.ParamSignUp) {
	// 1. 判断用户是否存在
	mysql.QueryUserByUsername()
	// 2. 生成uid
	snowflake.GenID()
	// 3. 保存到数据库
	mysql.InsertUser()
}
