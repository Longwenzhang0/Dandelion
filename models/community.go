package models

import "time"

type Community struct {
	ID   int64  `json:"id" db:"community_id"`     // 社区ID
	Name string `json:"name" db:"community_name"` // 社区名称
}

type CommunityDetail struct {
	ID           int64     `json:"id,string" db:"community_id"`              // 社区id
	Name         string    `json:"name" db:"community_name"`                 // 社区名称
	Introduction string    `json:"introduction,omitempty" db:"introduction"` // 社区简介，可省略
	CreateTime   time.Time `json:"create_time" db:"create_time"`             // 社区创建时间
}
