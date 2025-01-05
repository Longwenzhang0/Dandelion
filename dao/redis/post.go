package redis

import "Dandelion/models"

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
