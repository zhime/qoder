package middleware

import (
	"net/http"
	"strconv"
	"time"

	"devops/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	MaxRequests int           // 最大请求数
	Window      time.Duration // 时间窗口
	Message     string        // 限流消息
}

// DefaultRateLimitConfig 默认限流配置
var DefaultRateLimitConfig = RateLimitConfig{
	MaxRequests: 100,
	Window:      time.Minute,
	Message:     "请求过于频繁，请稍后再试",
}

// RateLimit 基于Redis的限流中间件
func RateLimit(rdb *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	cacheService := cache.NewCacheService(rdb, "devops")
	keys := cache.NewCacheKeys()

	return func(c *gin.Context) {
		// 获取用户ID，如果未登录则使用IP
		var identifier string
		if userID, exists := c.Get("user_id"); exists {
			identifier = "user:" + strconv.Itoa(int(userID.(uint)))
		} else {
			identifier = "ip:" + c.ClientIP()
		}

		// 获取API端点
		endpoint := c.Request.Method + ":" + c.FullPath()
		rateLimitKey := keys.APIRateLimit(0, identifier+":"+endpoint)

		ctx := c.Request.Context()

		// 获取当前请求计数
		count, err := cacheService.Increment(ctx, rateLimitKey)
		if err != nil {
			// Redis错误时允许请求通过
			c.Next()
			return
		}

		// 如果是第一次请求，设置过期时间
		if count == 1 {
			cacheService.SetExpire(ctx, rateLimitKey, config.Window)
		}

		// 检查是否超过限制
		if count > int64(config.MaxRequests) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": config.Message,
			})
			c.Abort()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(config.MaxRequests-int(count)))

		c.Next()
	}
}

// LoginRateLimit 登录限流中间件
func LoginRateLimit(rdb *redis.Client) gin.HandlerFunc {
	cacheService := cache.NewCacheService(rdb, "devops")
	keys := cache.NewCacheKeys()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		attemptKey := keys.LoginAttempts(ip)

		ctx := c.Request.Context()

		// 获取当前尝试次数
		attempts, err := cacheService.Increment(ctx, attemptKey)
		if err != nil {
			c.Next()
			return
		}

		// 第一次尝试，设置15分钟过期
		if attempts == 1 {
			cacheService.SetExpire(ctx, attemptKey, 15*time.Minute)
		}

		// 检查是否超过限制（5次）
		if attempts > 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "登录尝试次数过多，请15分钟后再试",
			})
			c.Abort()
			return
		}

		c.Next()

		// 如果登录成功，清除尝试计数
		if c.Writer.Status() == http.StatusOK {
			cacheService.Delete(ctx, attemptKey)
		}
	}
}
