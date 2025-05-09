package main_one

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	tokens     chan struct{} // Acts as the token bucket
	refillTime time.Duration // How often to add tokens
}

// NewRateLimiter initializes the rate limiter with capacity and refill interval
func NewRateLimiter(rateLimit int, refillTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(chan struct{}, rateLimit),
		refillTime: refillTime,
	}

	// Fill the bucket initially
	for i := 0; i < rateLimit; i++ {
		rl.tokens <- struct{}{}
	}

	go rl.startRefill()
	return rl
}

// startRefill periodically adds one token to the bucket, if it's not full
func (rl *RateLimiter) startRefill() {
	ticker := time.NewTicker(rl.refillTime)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case rl.tokens <- struct{}{}:
			// Token added
		default:
			// Bucket is full; skip
		}
	}
}

// Allow returns true if a token is available, false otherwise
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func main() {
	// Allow 5 requests every 10 seconds
	rateLimiter := NewRateLimiter(5, 10*time.Second)

	for i := 1; i <= 10; i++ {
		if rateLimiter.Allow() {
			fmt.Printf("[%2d] ✅ Request allowed\n", i)
		} else {
			fmt.Printf("[%2d] ❌ Request denied\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
