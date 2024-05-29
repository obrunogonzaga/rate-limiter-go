package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/obrunogonzaga/rate-limiter-go/pkg/limiter"
	"github.com/obrunogonzaga/rate-limiter-go/pkg/middleware"
)

func main() {
	godotenv.Load()

	redisUrl := os.Getenv("REDIS_URL")
	redisPort := os.Getenv("REDIS_PORT")
	newLimiter := limiter.NewLimiter(redisUrl, redisPort)

	r := gin.Default()
	r.Use(middleware.RateLimiterMiddleware(newLimiter))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the rate limited API",
		})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
