package logic

import (
	"Dandelion/dao/mysql"
	"Dandelion/models"
	"Dandelion/pkg/snowflake"
	_ "Dandelion/pkg/snowflake"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()

	// 2. 保存到数据库
	return mysql.CreatePost(p)
}

func GetPostByID(pid int64) (data *models.Post, err error) {
	return mysql.GetPostByID(pid)
}
