package limiter

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func setupTestLimiter(t *testing.T) (*Limiter, func()) {
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

	limiter := &Limiter{
		Client: rdb,
		Ctx:    context.Background(),
	}

	return limiter, func() {
		redisC.Terminate(ctx)
	}
}

func TestAllowRequest(t *testing.T) {
	limiter, teardown := setupTestLimiter(t)
	defer teardown()

	key := "test_ip"
	limit := 5
	blockTime := 10 * time.Second

	// Ensure Redis is clean before the test
	limiter.Client.Del(limiter.Ctx, key)

	for i := 0; i < limit; i++ {
		allowed := limiter.Allow(key, limit, blockTime)
		assert.True(t, allowed, "Request should be allowed")
	}

	// The next request should be blocked
	allowed := limiter.Allow(key, limit, blockTime)
	assert.False(t, allowed, "Request should be blocked after limit is reached")

	// Verify that the key:block is set in Redis
	blocked, err := limiter.Client.Get(limiter.Ctx, key+":block").Result()
	assert.Nil(t, err, "There should be a block key in Redis")
	assert.Equal(t, "blocked", blocked, "Block key should have the value 'blocked'")
}
