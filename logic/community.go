package logic

import (
	"Dandelion/dao/mysql"
	"Dandelion/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查询数据库，返回所有community
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
