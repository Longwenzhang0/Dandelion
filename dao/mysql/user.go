package mysql

import (
	"Dandelion/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// DAO层 把每一步数据库的操作都封装成函数，等待LOGIC层按需求调用

const secret = "longwenzhang0"

// CheckUserExist 查询指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
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
