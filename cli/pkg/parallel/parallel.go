package parallel

import (
	"context"
	"sync"
	"time"

	"github.com/epicchewy/converter/cli/pkg/logger"
)

type Executor interface {
	Do(ctx context.Context, i interface{}) (interface{}, error)
}

type Producer interface {
	Produce(ctx context.Context, queue Queue) error
}

type ResultHandler interface {
	Handle(ctx context.Context, i interface{}, err error)
}

type ParallelRunner struct {
	producer      Producer
	workers       []*Worker
	queue         Queue
	resultHandler ResultHandler
}

func NewRunner(producer Producer, executors []Executor, queue Queue, resultHandler ResultHandler) *ParallelRunner {
	workers := make([]*Worker, len(executors))
	for i, executor := range executors {
		workers[i] = NewWorker(executor, queue, resultHandler, i+1)
	}
	return &ParallelRunner{
		producer:      producer,
		workers:       workers,
		queue:         queue,
		resultHandler: resultHandler,
	}
}

func (r *ParallelRunner) Run(ctx context.Context) error {
	logger.Infow("Starting parallel runner", "num_workers", len(r.workers))

	var wg sync.WaitGroup
	for _, w := range r.workers {
		go w.Run(ctx, &wg)
	}

	defer func() {
		logger.Infow("Stopping workers and closing worker queue")
		r.queue.Stop()
		for _, w := range r.workers {
			w.Stop()
		}
		wg.Wait()
	}()

	if err := r.producer.Produce(ctx, r.queue); err != nil {
		logger.Errorw("Failed to produce work", "error", err)
		return err
	}

	isRunning := true
	for !r.queue.IsEmpty() && isRunning {
		select {
		case <-ctx.Done():
			isRunning = false
		default:
			logger.Infow("Waiting for queue to empty")
			time.Sleep(50 * time.Millisecond)
		}
	}

	return nil
}
