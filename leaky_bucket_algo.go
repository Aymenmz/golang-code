package main

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity int
	leakRate time.Duration
	tokens   int
	lastLeak time.Time
	mu       sync.Mutex
}

func NewLeakyBucket(capacity int, leakRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		tokens:   capacity,
		lastLeak: time.Now(),
	}
}

func (lb *LeakyBucket) allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elaspedTime := now.Sub(lb.lastLeak)
	tokenToAdd := int(elaspedTime / lb.leakRate)
	lb.tokens += tokenToAdd
	if lb.tokens > lb.capacity {
		lb.tokens = lb.capacity
	}

	lb.lastLeak = lb.lastLeak.Add(time.Duration(tokenToAdd) * lb.leakRate)

	fmt.Printf("Tokens added %d, Tokens substracted %d, Total tokens: %d\n", tokenToAdd, 1, lb.tokens)
	fmt.Printf("Laast leak time: %v\n", lb.lastLeak)
	if lb.tokens > 0 {
		lb.tokens--
		return true
	}
	return false

}

func main() {
	leakyBucketInstance := NewLeakyBucket(5, 500*time.Millisecond)
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if leakyBucketInstance.allow() {
				fmt.Printf("Current Time: %v\n", time.Now())
				fmt.Println("Request Allowed")
				fmt.Println("-----------------------------------------------------------------------------------------------------")
			} else {
				fmt.Printf("Current Time: %v\n", time.Now())
				fmt.Println("Request Denied")
				fmt.Println("-----------------------------------------------------------------------------------------------------")

			}
		}()
	}

	wg.Wait()
}
