package models

type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string // 数据库中没有该字段，为了返回整个user给前端
}
