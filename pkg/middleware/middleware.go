package middleware

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/obrunogonzaga/rate-limiter-go/pkg/limiter"
)

func RateLimiterMiddleware(limiter *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		godotenv.Load()

		limitIP := getEnvAsInt("RATE_LIMITER_IP_LIMIT", 5)
		limitToken := getEnvAsInt("LIMIT_PER_SECOND_TOKEN", 10)
		blockTime := time.Duration(getEnvAsInt("BLOCK_TIME", 300)) * time.Second

		ip := c.ClientIP()

		token := c.GetHeader("API_KEY")
		key := ip
		limit := limitIP

		if token != "" {
			key = token
			limit = limitToken
		}

		if limiter.IsBlocked(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
			return
		}

		if !limiter.Allow(key, limit, blockTime) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
			return
		}

		c.Next()
	}
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
