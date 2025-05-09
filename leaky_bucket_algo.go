package main

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity  int
	leakRate  time.Duration
	tokens    int
	lastLeak  time.Time
	mu        sync.Mutex
}

// NewLeakyBucket creates a new leaky bucket rate limiter
func NewLeakyBucket(capacity int, leakRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		tokens:   capacity,
		lastLeak: time.Now(),
	}
}

// Allow checks if a request can be processed based on the leak rate
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsedTime := now.Sub(lb.lastLeak)

	// Calculate how many tokens to leak (i.e., how many slots freed)
	tokensToAdd := int(elapsedTime / lb.leakRate)
	if tokensToAdd > 0 {
		lb.tokens += tokensToAdd
		if lb.tokens > lb.capacity {
			lb.tokens = lb.capacity
		}
		lb.lastLeak = lb.lastLeak.Add(time.Duration(tokensToAdd) * lb.leakRate)
	}

	fmt.Printf("⏳ Tokens added: %d | Total tokens: %d | Time: %v\n", tokensToAdd, lb.tokens, lb.lastLeak)

	if lb.tokens > 0 {
		lb.tokens--
		return true
	}
	return false
}

func main() {
	bucket := NewLeakyBucket(5, 500*time.Millisecond)
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if bucket.Allow() {
				fmt.Printf("[%2d] ✅ Request allowed at %v\n", id, time.Now())
			} else {
				fmt.Printf("[%2d] ❌ Request denied at %v\n", id, time.Now())
			}
			fmt.Println("------------------------------------------------")
		}(i)
		time.Sleep(200 * time.Millisecond) // simulate staggered requests
	}

	wg.Wait()
}
