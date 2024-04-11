package ratelimiter

import (
	"sync"

	"go.uber.org/ratelimit"
)

type RateLimiterBuilder struct {
	Key       string
	Threshold int
}

type SyncedLimiter struct {
	Limiter ratelimit.Limiter
	mux     sync.Mutex
}

// RateLimiterBasedOnKey is a struct that holds a map of ratelimit.Limiter
type RateLimiterBasedOnKey struct {
	limiters map[string]*SyncedLimiter
}

func (builder *RateLimiterBuilder) build() *RateLimiterBasedOnKey {
	syncedLimiter := &SyncedLimiter{
		Limiter: ratelimit.New(builder.Threshold),
	}
	limiters := &RateLimiterBasedOnKey{
		limiters: make(map[string]*SyncedLimiter), // Change the type of map value to *ratelimit.Limiter
	}

	limiters.limiters[builder.Key] = syncedLimiter

	return limiters
}

func NewRateLimiter(Key string, threshold int) *RateLimiterBasedOnKey {
	builder := &RateLimiterBuilder{
		Key:       Key,
		Threshold: threshold,
	}

	return builder.build()
}

// Take is a method to take a token from the rate limiter
func (r *RateLimiterBasedOnKey) Take(key string) {
	limiter := r.limiters[key]
	limiter.mux.Lock()
	defer limiter.mux.Unlock()
	limiter.Limiter.Take()
}
