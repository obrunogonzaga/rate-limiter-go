package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/obrunogonzaga/rate-limiter-go/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}

	ip, err := redisC.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %s", err)
	}

	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}

	redisAddr := ip + ":" + port.Port()
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return rdb, func() {
		redisC.Terminate(ctx)
	}
}

func setEnvVariables() {
	os.Setenv("LIMIT_PER_SECOND_IP", "5")
	os.Setenv("LIMIT_PER_SECOND_TOKEN", "10")
	os.Setenv("BLOCK_TIME_SECONDS", "300")
}

func TestRateLimiterMiddlewareByIP(t *testing.T) {
	rdb, teardown := setupTestRedis(t)
	defer teardown()

	setEnvVariables()

	l := &limiter.Limiter{
		Client: rdb,
		Ctx:    context.Background(),
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimiterMiddleware(l))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Simulate requests from the same IP
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// The next request should be blocked
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimiterMiddlewareByToken(t *testing.T) {
	rdb, teardown := setupTestRedis(t)
	defer teardown()

	setEnvVariables()

	l := &limiter.Limiter{
		Client: rdb,
		Ctx:    context.Background(),
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimiterMiddleware(l))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Simulate requests with the same token
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("API_KEY", "abc123")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// The next request should be blocked
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "abc123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
