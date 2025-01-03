package redis

// redis key，注意使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix              = "dandelion:"
	KeyPostTime            = "post:time"  // zset;帖子及发帖时间
	KeyPostScore           = "post:score" // zset;帖子及分数
	KeyPostVotedZSetPrefix = "post:voted" // zset;（非完整key）记录用户及投票类型；参数为post id
)
