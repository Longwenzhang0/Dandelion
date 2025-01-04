package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"

	"go.uber.org/zap"
)

/* 投票功能的多种情况梳理,为帖子投票的函数，投一票就加432分；一天86400s/200，200张赞成票就可以把帖子续一天
direction = 1,
	1. 之前没有投过票，现在投赞成 	-->更新分数和投票记录，diff=1，+432
	2. 之前投反对，现在改成赞成		-->更新分数和投票记录，diff=2，+432*2
direction = 0
	1. 之前投反对，现在取消			-->更新分数和投票记录，diff=1，+432		-->以上三条，现在的direction大于之前的；以下都是小于
	2. 之前投赞成，现在取消			-->更新分数和投票记录，diff=1，-432
direction = -1
	1. 之前没有投过票，现在投反对票	-->更新分数和投票记录，diff=1，-432
	2. 之前投赞成，现在改成反对		-->更新分数和投票记录，diff=2，-432*2
投票的限制：
	每个帖子发表一周内允许投票，超过一个周不允许投票；
	1. 到期之后将redis中保存的赞成票数or反对票 存储到mysql
	2. 到期之后删除key：KeyPostVotedZSetPrefix
*/

const (
	oneWeekSeconds = 7 * 24 * 3600 // 一周的秒数，用于判定超时
	scorePerVote   = 432           // 每票的分数
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID int64) error {
	// 以下两个操作需要同时成功，使用Pipeline
	pipeline := client.TxPipeline()

	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 执行
	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("CreatePost/pipeline.Exec() failed: ", zap.Error(err))
	}
	return err
}

func VoteForPost(userID, postID string, direction float64) error {
	// 1. 判断投票限制
	// 获取帖子发布时间，超时返错并记录日志
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekSeconds {
		zap.L().Error("VoteForPost/float64(time.Now().Unix())-postTime > oneWeekSeconds", zap.Error(ErrorVoteTimeExpire))
		return ErrorVoteTimeExpire
	}

	// 2和3需要在同一个pipeline里操作
	pipeline := client.TxPipeline()

	// 2. 更新帖子的分数

	// 查询当前用户给当前帖子之前的投票记录
	oldDirection := client.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	// 判断是加分还是减分，diff为投票数的绝对值，将多个分支转换成一个公式
	var operator float64
	if direction > oldDirection {
		operator = 1
	} else {
		operator = -1
	}
	diff := math.Abs(oldDirection - direction)
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), operator*diff*scorePerVote, postID)

	// 3. 记录用户为该帖子投票的数据
	if direction == 0 {
		// 移除
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  direction, // 当前用户投的是赞成还是反对
			Member: userID,
		})
	}
	// 执行
	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("VoteForPost/pipeline.Exec() failed: ", zap.Error(err))
	}
	return err
}
