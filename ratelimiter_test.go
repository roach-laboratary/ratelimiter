package ratelimiter_test

import (
	"sync"
	"testing"
	"time"

	"github.com/roach-laboratary/ratelimiter"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterBuilder_build(t *testing.T) {
	key := "test"
	limiter := ratelimiter.NewRateLimiter(key, 1)

	assert.NotNil(t, limiter)
}

func TestRateLimiter_Take(t *testing.T) {
	key := "test"
	limiter := ratelimiter.NewRateLimiter(key, 1)

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Call Take method
			limiter.Take(key)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	if elapsed < 9*time.Second {
		t.Errorf("Execution time was %v, expected it to be over 9 seconds", elapsed)
	}
}

func TestRateLimiter_Tak2(t *testing.T) {
	key1 := "test"
	key2 := "test2"
	limiter := ratelimiter.NewRateLimiter(key1, 1)
	limiter2 := ratelimiter.NewRateLimiter(key2, 1)

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Call Take method
			limiter.Take(key1)
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Call Take method
			limiter2.Take(key2)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	if elapsed > 10*time.Second {
		t.Errorf("Execution time was %v, expected it to be under 10 seconds", elapsed)
	}
}

func TestTakeIfThereisNoKey(t *testing.T) {
	key := "test"
	limiter := ratelimiter.NewRateLimiter(key, 1)

	err := limiter.Take("test2")
	assert.NotNil(t, err)
}

func TestCreateKeyIfThereisNoKey(t *testing.T) {
	key := "test"
	limiter := ratelimiter.NewRateLimiter(key, 1)
	limiter.CreateKey("test2", 1)

	err := limiter.Take("test2")
	assert.Nil(t, err)
}

func TestCreateKeyIfThereisAlreadyKey2(t *testing.T) {
	key := "test"
	limiter := ratelimiter.NewRateLimiter(key, 1)
	limiter.CreateKey("test2", 1)
	limiter.CreateKey("test2", 1)

	err := limiter.Take("test2")
	assert.Nil(t, err)
}
