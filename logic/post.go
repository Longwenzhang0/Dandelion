package logic

import (
	"Dandelion/dao/mysql"
	"Dandelion/dao/redis"
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
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return
}

// GetPostByID 根据单个id查询帖子详情
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

// GetPostList 分页查询帖子详情
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 获取到所有帖子信息
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("GetPostList/mysql.GetPostList() failed: ", zap.Error(err))
		return nil, err
	}

	// 初始化data，根据帖子数量来决定data的大小
	data = make([]*models.ApiPostDetail, 0, len(posts))

	// 遍历帖子中的用户信息和社区信息
	for _, post := range posts {
		// 根据帖子信息中的AuthorID查询用户信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetPostList/mysql.GetUserByID(post.AuthorID) failed: ", zap.Error(err))
			continue
		}
		// 根据帖子信息中的CommunityID查询社区信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetPostList/mysql.GetCommunityDetailByID(post.CommunityID) failed: ", zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("GetPostList2/redis.GetPostIDsInOrder(p) failed:", zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostList2/redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 3. 根据id去mysql查询帖子的详细信息
	posts, err := mysql.GetPostByIDs(ids)
	if err != nil {
		zap.L().Error("GetPostList2/mysql.GetPostByIDs(ids) failed: ", zap.Error(err))
		return
	}
	// 初始化data，根据帖子数量来决定data的大小
	data = make([]*models.ApiPostDetail, 0, len(posts))

	// 提前查询每篇帖子的投票数量，跟帖子数量是一一对应的，直接填充即可。
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("GetPostList2/redis.GetPostVoteData(ids) failed: ", zap.Error(err))
		return
	}

	// 遍历帖子中的用户信息和社区信息
	for idx, post := range posts {
		// 根据帖子信息中的AuthorID查询用户信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetPostList/mysql.GetUserByID(post.AuthorID) failed: ", zap.Error(err))
			continue
		}
		// 根据帖子信息中的CommunityID查询社区信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetPostList/mysql.GetCommunityDetailByID(post.CommunityID) failed: ", zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}
