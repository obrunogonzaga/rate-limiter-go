package limiter

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	ctx    context.Context
}

func NewLimiter(redisUrl, redisPort string) *Limiter {
	rbr := redis.NewClient(&redis.Options{
		Addr: redisUrl + redisPort,
	})

	return &Limiter{
		client: rbr,
		ctx:    context.Background(),
	}
}

func (l *Limiter) Allow(key string, limit int, blockTime time.Duration) bool {
	current, err := l.client.Get(l.ctx, key).Int()
	if err != nil {
		return l.handleGetError(err, key, blockTime)
	}
	return l.handleCurrentLimit(current, key, limit, blockTime)
}

func (l *Limiter) IsBlocked(key string) bool {
	_, err := l.client.Get(l.ctx, key+":block").Result()
	return !errors.Is(err, redis.Nil)
}

func (l *Limiter) handleGetError(err error, key string, blockTime time.Duration) bool {
	if errors.Is(err, redis.Nil) {
		return l.incrementKeyAndSetExpiry(key, blockTime)
	}
	log.Printf("Error getting key %s: %v", key, err)
	return false
}

func (l *Limiter) handleCurrentLimit(current int, key string, limit int, blockTime time.Duration) bool {
	if current < limit {
		return l.incrementKeyAndSetExpiry(key, blockTime)
	}
	if current >= limit {
		l.client.Set(l.ctx, key+":block", "blocked", blockTime)
		return false
	}
	return true
}

func (l *Limiter) incrementKeyAndSetExpiry(key string, blockTime time.Duration) bool {
	pipeline := l.client.TxPipeline()
	pipeline.Incr(l.ctx, key)
	pipeline.Expire(l.ctx, key, blockTime)
	_, err := pipeline.Exec(l.ctx)
	if err != nil {
		log.Printf("Error executing pipeline for key %s: %v", key, err)
		return false
	}
	return true
}
