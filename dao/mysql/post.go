package mysql

import (
	"Dandelion/models"
	"database/sql"

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

// 根据post_id获取帖子详情
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
