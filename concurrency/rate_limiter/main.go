package main

import (
	"context"
	"fmt"
	"time"
)

type RateLimiter struct {
	ticker   *time.Ticker
	tokens   chan struct{}
	stopChan chan struct{}
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
	if rate <= 0 {
		rate = 1
	}

	rl := &RateLimiter{
		ticker:   time.NewTicker(time.Second / time.Duration(rate)),
		tokens:   make(chan struct{}, burst),
		stopChan: make(chan struct{}),
	}

	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-rl.ticker.C:
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
			case <-rl.stopChan:
				rl.ticker.Stop()
				fmt.Println("Stop")
				return
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.tokens:
		return nil
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

func main() {
	rl := NewRateLimiter(3, 5)
	defer rl.Stop()

	ctx := context.Background()
	for i := 0; i < 20; i++ {
		if err := rl.Wait(ctx); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Request", i)
	}
}
