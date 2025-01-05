package models

import "time"

// Post 帖子结构体，相同类型的数据放在一起，会更省空间；内存对齐
type Post struct {
	ID          int64     `json:"id" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id"  binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title"  binding:"required"`
	Content     string    `json:"content" db:"content"  binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详情接口的结构体；嵌入帖子结构体和社区结构体
type ApiPostDetail struct {
	AuthorName       string `json:"author_name"`
	VoteNum          int64  `json:"vote_num"`
	*Post            `json:"community_post"`
	*CommunityDetail `json:"community_detail"`
}
