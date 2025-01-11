package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// RateLimitMiddleware 限流中间件，采用令牌桶算法
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌，abort；取1个令牌，没有的时候会返回0；或者bucket.Take(1) >0
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit")
			c.Abort()
			return
		}
		// 反之放行next
		c.Next()
	}
}
