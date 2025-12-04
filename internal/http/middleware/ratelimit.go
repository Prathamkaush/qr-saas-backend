package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func RateLimit(redis *redis.Client, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		key := "rl:" + route + ":" + ip
		pipe := redis.TxPipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, err := pipe.Exec(ctx)
		if err != nil {
			// on redis error, just continue
			c.Next()
			return
		}

		count := incr.Val()
		if int(count) > maxRequests {
			c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(maxRequests-int(count)))
		c.Next()
	}
}
