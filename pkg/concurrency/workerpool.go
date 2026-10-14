package concurrency

import (
	"context"
	"sync"
)

// Task represents a unit of work executed by the pool
type Task func()

// Pool manages a pool of worker goroutines utilizing the GMP scheduler
type Pool struct {
	tasksChan chan Task
	wg        sync.WaitGroup
	workers   int
}

// NewPool creates an initialized worker pool
func NewPool(workers int, queueSize int) *Pool {
	return &Pool{
		tasksChan: make(chan Task, queueSize),
		workers:   workers,
	}
}

// Start launches worker goroutines that process tasks until context is cancelled
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-p.tasksChan:
					if !ok {
						return
					}
					task()
				}
			}
		}(i)
	}
}

// Submit queues a task to be executed by the pool
func (p *Pool) Submit(task Task) {
	p.tasksChan <- task
}

// Stop closes the task queue and waits for all active workers to terminate
func (p *Pool) Stop() {
	close(p.tasksChan)
	p.wg.Wait()
}
