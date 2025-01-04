package logic

import (
	"Dandelion/dao/redis"
	"Dandelion/models"
	"strconv"

	"go.uber.org/zap"
)

// VoteForPost 为帖子投票的函数，投一票就加432分；一天86400s/200，200张赞成票就可以把帖子续一天
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.Int64("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.FormatInt(userID, 10), strconv.FormatInt(p.PostID, 10), float64(p.Direction))
}
