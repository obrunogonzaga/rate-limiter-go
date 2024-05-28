package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("RateLimiterMiddleware")
		c.Next()
	}
}
