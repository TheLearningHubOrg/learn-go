package concurrency

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := NewPool(4, 100)
	pool.Start(ctx)

	var counter int64
	totalTasks := 100

	for i := 0; i < totalTasks; i++ {
		pool.Submit(func() {
			atomic.AddInt64(&counter, 1)
		})
	}

	// Allow execution
	time.Sleep(50 * time.Millisecond)
	pool.Stop()

	if atomic.LoadInt64(&counter) != int64(totalTasks) {
		t.Fatalf("expected %d tasks to complete, got %d", totalTasks, counter)
	}
}
