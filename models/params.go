package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 定义请求的参数结构体

// ParamSignUp 注册参数结构体；可以在原有基础上新增tag
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// ParamLogin 登录结构体；
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据结构体，谁给哪个帖子投票，需要UserID和PostID，UserID可以从token中解析出来
type ParamVoteData struct {
	//UserID	// 从请求中获取当前用户
	PostID    int64 `json:"post_id,string" binding:"required"`       // 帖子id
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成1，反对-1，取消投票0；oneof=1 0 -1表示只能是其中一个的值
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}
