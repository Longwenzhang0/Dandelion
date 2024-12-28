package mysql

import (
	"Dandelion/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// DAO层 把每一步数据库的操作都封装成函数，等待LOGIC层按需求调用

const secret = "longwenzhang0"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
)

// CheckUserExist 查询指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 插入用户
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密，更新回原结构体
	user.Password = encryptPassword(user.Password)

	// 执行sql语句，存入数据库
	sqlStr := `insert into user(user_id,username,password) values (?,?,?)`
	db.Exec(sqlStr, user.UserID, user.Username, user.Password)

	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	sum := h.Sum([]byte(oPassword))

	return hex.EncodeToString(sum)
}

// Login 登录
func Login(user *models.User) (err error) {
	// 保存明文密码
	oPassword := user.Password
	sqlStr := `select user_id,username,password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username)
	// 两种错误，一种是用户不存在，一种是查询失败
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	// 如果查询成功，则判断返回的密码；
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}
