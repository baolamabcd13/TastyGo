package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Lưu trữ các rate limiter cho từng IP
var (
    ipLimiters = make(map[string]*rate.Limiter)
    mu         sync.Mutex
)

// Tạo rate limiter cho IP
func getIPLimiter(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()

    limiter, exists := ipLimiters[ip]
    if !exists {
        // Cho phép 5 request mỗi phút
        limiter = rate.NewLimiter(rate.Every(time.Minute/5), 5)
        ipLimiters[ip] = limiter
    }

    return limiter
}

// Middleware giới hạn tần suất request theo IP
func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        limiter := getIPLimiter(ip)
        
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
                "message": "too many requests, please try again later",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Dọn dẹp rate limiters không sử dụng
func CleanupRateLimiters() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        mu.Lock()
        for ip, limiter := range ipLimiters {
            // Xóa limiter nếu không có request trong 1 giờ
            if limiter.Tokens() >= float64(limiter.Burst()) {
                delete(ipLimiters, ip)
            }
        }
        mu.Unlock()
    }
}
