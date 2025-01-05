package redis

import (
	"Dandelion/models"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// GetPostIDsInOrder 从redis中查找符合规则的post id，降序
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 根据请求的排序参数来指定redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 确定查询起始编号
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 根据key,start,end从redis中查询，分数从大到小
	return client.ZRevRange(key, start, end).Result()
}

// GetPostVoteData 根据ids查询每篇帖子的赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// 初始化pipeline和结果slice；批量发送命令，减少rtt
	pipeline := client.Pipeline()
	data = make([]int64, 0, len(ids))

	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		// 先将命令放入pipeline，等待统一执行
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("GetPostVoteData/pipeline.Exec() failed: ", zap.Error(err))
		return nil, err
	}
	// 写入结果slice中
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
