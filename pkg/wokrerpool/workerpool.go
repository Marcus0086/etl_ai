package wokrerpool

import (
	"context"
	"sync"
)

type WorkerPool[T any] struct {
	numWorkers int
	jobQueue   chan T
	wg         sync.WaitGroup
	errors     []error
	mu         sync.Mutex
}

type WorkerFunc[T any] func(workerID int, chunk T) error

func New[T any](numWorkers int) *WorkerPool[T] {
	return &WorkerPool[T]{
		numWorkers: numWorkers,
		jobQueue:   make(chan T),
	}
}

func (wp *WorkerPool[T]) Start(ctx context.Context, workerFunc WorkerFunc[T]) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i, workerFunc)
	}
}

func (wp *WorkerPool[T]) Submit(job T) {
	wp.jobQueue <- job
}

func (wp *WorkerPool[T]) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
}

func (wp *WorkerPool[T]) worker(ctx context.Context, workerID int, workerFunc WorkerFunc[T]) {
	defer wp.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			err := workerFunc(workerID, job)
			if err != nil {
				wp.mu.Lock()
				defer wp.mu.Unlock()
				wp.errors = append(wp.errors, err)
			}
		}
	}
}

func (wp *WorkerPool[T]) Errors() []error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.errors
}
