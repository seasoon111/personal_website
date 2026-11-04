package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"personal-website/backend/internal/pkg/response"
)

// RateLimit 按 IP 限流：window 时间内最多 limit 次请求。
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 5
	}
	if window <= 0 {
		window = time.Minute
	}
	interval := window / time.Duration(limit)

	type entry struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		entries = make(map[string]*entry)
	)

	// 定期清理 10 分钟未活跃的 IP，避免内存无限增长。
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			for ip, e := range entries {
				if time.Since(e.lastSeen) > 10*time.Minute {
					delete(entries, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		e, ok := entries[ip]
		if !ok {
			e = &entry{limiter: rate.NewLimiter(rate.Every(interval), limit)}
			entries[ip] = e
		}
		e.lastSeen = time.Now()
		mu.Unlock()

		if !e.limiter.Allow() {
			response.Fail(c, http.StatusTooManyRequests, 42900, "操作过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

// LoginLimiter 登录接口限流（默认 5 次/分钟/IP）。
func LoginLimiter(perMinute int) gin.HandlerFunc {
	return RateLimit(perMinute, time.Minute)
}
