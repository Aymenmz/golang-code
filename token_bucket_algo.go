package main_one

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	tokens     chan struct{} // token storage (channel buffer acts like the bucket)
	refillTime time.Duration // time interval between each token refill
}

func NewRateLimiter(rateLimite int, refillTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(chan struct{}, rateLimite),
		refillTime: refillTime,
	}
	for range rateLimite {
		rl.tokens <- struct{}{}
	}
	go rl.startRefill()
	return rl
}

func (rt *RateLimiter) startRefill() {
	ticker := time.NewTicker(rt.refillTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			select {
			case rt.tokens <- struct{}{}:
			default:
			}
		}
	}
}

func (rt *RateLimiter) allow() bool {
	select {
	case <-rt.tokens:
		return true
	default:
		return false

	}
}

func main() {
	rateLimiter := NewRateLimiter(5, 10*time.Second) 
	//(numRequest, duration of refill) so each refill duration, we allow new numRequest credits 
	// 100 200 300 400 500 | 
	for range 10 {
		if rateLimiter.allow() {
			fmt.Println("Request allowed")
		} else {
			fmt.Println("Request denied")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
