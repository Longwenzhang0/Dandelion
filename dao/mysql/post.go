package mysql

import (
	"Dandelion/models"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post
	(post_id, title, content, author_id, community_id)
	values(?, ?, ?, ?, ?)
	`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	if err != nil {
		zap.L().Error("CreatePost/db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID) failed: ", zap.Error(err))
	}
	return
}

// GetPostByID 根据post_id获取帖子详情
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
    post_id,title,content,author_id,community_id,create_time 
	from post
	where post_id = ?
	`
	err = db.Get(post, sqlStr, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
			zap.L().Error("GetPostByID/db.Get(post, sqlStr, pid) failed: ", zap.Error(err))
		}
	}
	return post, err

}

// GetPostByIDs 根据id列表查询帖子详情
func GetPostByIDs(ids []string) (postList []*models.Post, err error) {
	// 查询时指定post_id以传入的id顺序来排序
	sqlStr := `select
	post_id,title,content,author_id,community_id,create_time
	from post 
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)
	`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		zap.L().Error("GetPostByIDs/sqlx.In() failed: ", zap.Error(err))
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

// 返回所有帖子的slice
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	// 每次限制两条
	sqlStr := `select
	post_id,title,content,author_id,community_id,create_time
	from post 
	order by create_time
	desc 
	limit ?,?
	`
	// limit两个参数，分别是偏移量和行数。填充时计算当前的偏移量为(page-1)*size，比如要取1页的10条，那起始位置就是0，offset为10
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
			zap.L().Error("GetPostList/db.Select(&posts, sqlStr) failed: ", zap.Error(err))
		}
	}
	return
}
