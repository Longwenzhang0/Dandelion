package logic

import (
	"Dandelion/dao/mysql"
	"Dandelion/models"
	"Dandelion/pkg/snowflake"
	_ "Dandelion/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()

	// 2. 保存到数据库
	return mysql.CreatePost(p)
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合接口想用的数据
	// 查询帖子信息
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostByID/mysql.GetPostByID(pid) failed: ", zap.Error(err))
		return
	}
	// 根据帖子信息中的AuthorID查询用户信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("GetPostByID/mysql.GetUserByID(post.AuthorID) failed: ", zap.Error(err))
		return
	}
	// 根据帖子信息中的CommunityID查询社区信息
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("GetPostByID/mysql.GetCommunityDetailByID(post.CommunityID) failed: ", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}

	return
}
