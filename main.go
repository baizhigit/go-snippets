package main

import (
	"context"
	"fmt"
	"time"

	cc "github.com/baizhigit/go-snippets/concurrency"
)

func main() {
	println("===Snipetts Main run===")
	concurrency()
}

func concurrency() {
	println("===Snipetts Concurrency run===")

	// Worker Pool Pattern
	println("1. Worker Pool Pattern run")
	const (
		numWorkers = 3
		numTasks   = 10
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create worker pool
	pool := cc.NewWorkerPool(numWorkers, numTasks, cc.ProcessTask)

	// Start workers
	pool.Start(ctx)

	// Submit tasks in a separate goroutine
	go func() {
		for i := 0; i < numTasks; i++ {
			task := cc.Task{
				ID:   i,
				Name: fmt.Sprintf("Task_%d.exe", i),
			}
			if err := pool.Submit(task); err != nil {
				fmt.Printf("❌ Failed to submit task %d: %v\n", i, err)
			} else {
				fmt.Printf("📥 Submitted task %d\n", i)
			}
		}

		pool.Shutdown()
	}()

	// Display result
	fmt.Println("\n=== Results ===")
	for res := range pool.Results() {
		fmt.Printf("✅ %s\n", res)
	}
}
