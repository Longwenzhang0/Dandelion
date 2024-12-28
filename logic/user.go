package logic

import (
	"Dandelion/dao/mysql"
	"Dandelion/models"
	"Dandelion/pkg/snowflake"
)

// LOGIC层 用于存放业务逻辑的代码

// SignUp 注册用户
func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在；返回值是err，包含查询出错和已存在的两种情况
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2. 生成uid
	userID := snowflake.GenID()
	// 构造一个Users实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3. 保存到数据库
	err = mysql.InsertUser(user)
	return
}

// Login 用户登录
func Login(p *models.ParamLogin) (err error) {
	// 构造User实例
	user := &models.User{
		UserID:   0,
		Username: p.Username,
		Password: p.Password,
	}
	err = mysql.Login(user)
	return
}
