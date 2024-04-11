package ratelimiter

import (
	"errors"
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
	mux      sync.Mutex
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
func (r *RateLimiterBasedOnKey) Take(key string) error {
	if limiter, exists := r.limiters[key]; exists {
		limiter.mux.Lock()
		defer limiter.mux.Unlock()
		limiter.Limiter.Take()
		return nil
	}

	// throw an error if the key does not exist
	return errors.New("key does not exist")
}

// Create key is a method to create a new key in the rate limiter
func (r *RateLimiterBasedOnKey) CreateKey(key string, threshold int) {
	r.mux.Lock()
	defer r.mux.Unlock()
	syncedLimiter := &SyncedLimiter{
		Limiter: ratelimit.New(threshold),
	}
	r.limiters[key] = syncedLimiter
}
