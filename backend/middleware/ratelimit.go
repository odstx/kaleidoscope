package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redis             *redis.Client
	requestsPerMinute int
}

func NewRateLimiter(redis *redis.Client, requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		redis:             redis,
		requestsPerMinute: requestsPerMinute,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		var key string

		userID, exists := c.Get("userID")
		if exists {
			key = fmt.Sprintf("rate_limit:user:%v", userID)
		} else {
			ip := c.ClientIP()
			key = fmt.Sprintf("rate_limit:ip:%s", ip)
		}

		allowed, err := rl.checkRateLimit(ctx, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string) (bool, error) {
	now := time.Now()
	windowStart := now.Truncate(time.Minute)

	countKey := fmt.Sprintf("%s:%d", key, windowStart.Unix())

	count, err := rl.redis.Incr(ctx, countKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rl.redis.Expire(ctx, countKey, time.Minute)
	}

	if int(count) > rl.requestsPerMinute {
		return false, nil
	}

	return true, nil
}
