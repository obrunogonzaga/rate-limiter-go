package main

import (
	"github.com/gin-gonic/gin"
	"github.com/obrunogonzaga/rate-limiter-go/pkg/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.RateLimiterMiddleware())

	r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to the rate limited API",
        })
    })

	r.Run(":8080")
}
