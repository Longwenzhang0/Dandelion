package redis

// redis key，注意使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix              = "dandelion:"
	KeyPostTimeZSet        = "post:time"  // zset;帖子id及发帖时间
	KeyPostScoreZSet       = "post:score" // zset;帖子id及分数
	KeyPostVotedZSetPrefix = "post:voted" // zset;（非完整key）记录用户及投票类型；参数为post id
	KeyCommunitySetPrefix  = "community:" // set; 保存每个社区下所有帖子的id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
