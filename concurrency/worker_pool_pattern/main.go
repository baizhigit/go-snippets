package main

import (
	"fmt"
	"sync"
	"time"
)

func workerPool(jobs <-chan int, results chan<- int, numWorkers int) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- processJob(job)
			}
		}()
	}

	wg.Wait()
	close(results)
}

func processJob(num int) int {
	time.Sleep(time.Millisecond * 800)
	return num * 2
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	go workerPool(jobs, results, 5)

	for i := 0; i < 100; i++ {
		jobs <- i
	}
	close(jobs)

	for result := range results {
		fmt.Println(result)
	}
}
