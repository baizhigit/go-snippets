package concurrency

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Worker Pool Pattern
// Ограниченное количество горутин обрабатывающих задачи

type Task struct {
	ID   int
	Name string
}

type TaskProcessor func(Task) string

type WorkerPool struct {
	numWorkers int
	taskCh     chan Task
	resCh      chan string
	wg         sync.WaitGroup
	processor  TaskProcessor
}

func NewWorkerPool(numWorkers int, numTasks int, processor TaskProcessor) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		taskCh:     make(chan Task, numTasks),
		resCh:      make(chan string, numTasks),
		processor:  processor,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.Worker(ctx, i)
	}
}

func (wp *WorkerPool) Worker(ctx context.Context, workerID int) {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.taskCh:
			if !ok {
				fmt.Printf("Worker %d: shutting down (task channel closed)\n", workerID)
				return
			}

			fmt.Printf("Worker %d: started task %d (%s)\n", workerID, task.ID, task.Name)
			result := wp.processor(task)

			select {
			case wp.resCh <- result:
				fmt.Printf("Worker %d: completed task %d\n", workerID, task.ID)
			case <-ctx.Done():
				fmt.Printf("Worker %d: context cancelled sending result for task %d\n", workerID, task.ID)
				return
			}

		case <-ctx.Done():
			fmt.Printf("Worker %d: shutting down (context cancelled)\n", workerID)
			return
		}
	}
}

func (wp *WorkerPool) Submit(task Task) error {
	select {
	case wp.taskCh <- task:
		return nil
	default:
		return fmt.Errorf("task channel full, cannot submit task %d", task.ID)
	}
}

func (wp *WorkerPool) Shutdown() {
	close(wp.taskCh)
	wp.wg.Wait()
	close(wp.resCh)
}

func (wp *WorkerPool) Results() <-chan string {
	return wp.resCh
}

func ProcessTask(task Task) string {
	processingTime := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(processingTime)
	return fmt.Sprintf("Task %s processed successfully (ID %d) (took %v)", task.Name, task.ID, processingTime)
}
