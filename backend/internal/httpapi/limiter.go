package httpapi

import (
	"sync"
	"time"
)

type attemptWindow struct {
	count     int
	expiresAt time.Time
}

type loginLimiter struct {
	mu       sync.Mutex
	attempts map[string]attemptWindow
	limit    int
	window   time.Duration
	now      func() time.Time
}

func newLoginLimiter(limit int, window time.Duration) *loginLimiter {
	return &loginLimiter{
		attempts: make(map[string]attemptWindow),
		limit:    limit,
		window:   window,
		now:      time.Now,
	}
}

func (limiter *loginLimiter) allow(key string) bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := limiter.now()
	attempt := limiter.attempts[key]
	if now.After(attempt.expiresAt) {
		attempt = attemptWindow{expiresAt: now.Add(limiter.window)}
	}
	if attempt.count >= limiter.limit {
		limiter.attempts[key] = attempt
		return false
	}
	attempt.count++
	limiter.attempts[key] = attempt
	return true
}

func (limiter *loginLimiter) reset(key string) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	delete(limiter.attempts, key)
}
