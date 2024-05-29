package limiter

type RateLimiterStrategy interface {
	Allow() bool
	IsBlocked() bool
}