package redis

import (
	"Dandelion/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 确定查询起始编号
	start := (page - 1) * size
	end := start + size - 1
	// 根据key,start,end从redis中查询，分数从大到小
	return client.ZRevRange(key, start, end).Result()
}

// GetPostIDsInOrder 从redis中查找符合规则的post id，降序
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 根据请求的排序参数来指定redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
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

// GetCommunityPostIDsInOrder 按照社区，从redis中查找符合规则的post id，降序
func GetCommunityPostIDsInOrder(p *models.ParamCommunityPostList) ([]string, error) {
	// 根据请求的排序参数来指定redis key
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用zinterstore把社区的帖子set，与帖子分数的zset 生成一个zset
	// 在新生成的zset中，按照之前的逻辑查询数据

	// 初始化社区Key
	communityIDStr := strconv.FormatInt(p.CommunityID, 10) // 将int64转为str
	communityKey := getRedisKey(KeyCommunitySetPrefix + communityIDStr)
	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + communityIDStr
	if client.Exists(key).Val() < 1 {
		// 不存在该key，需要计算；参数为MAX
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, communityKey, orderKey) // Zinterstore操作

		pipeline.Expire(key, 24*time.Hour) // 设置超时时间

		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的话直接根据key查询ids
	return getIDsFromKey(key, p.Page, p.Size)
}
