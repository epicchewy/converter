package parallel

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/epicchewy/converter/cli/pkg/logger"
)

type Worker struct {
	executor       Executor
	queue          Queue
	resultHandler  ResultHandler
	workerQuitChan chan bool
	id             int
	running        atomic.Bool // Only used for testing
}

func NewWorker(executor Executor, queue Queue, resultHandler ResultHandler, id int) *Worker {
	return &Worker{
		executor:       executor,
		queue:          queue,
		resultHandler:  resultHandler,
		workerQuitChan: make(chan bool),
		id:             id,
	}
}

func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	w.running.Store(true)
	defer func() { w.running.Store(false) }()

	for {
		select {
		case <-w.workerQuitChan:
			logger.Infof("Worker %d received quit signal, shutting down", w.id)
			return
		case <-ctx.Done():
			logger.Infof("Worker %d received context done signal, shutting down", w.id)
			return
		default:
		}

		work, err := w.queue.Dequeue(ctx)
		if err != nil {
			logger.Infof("Got nil work, shutting down worker %d", w.id)
			return
		}

		result, err := w.executor.Do(ctx, work)
		w.resultHandler.Handle(ctx, result, err)
	}
}

func (w *Worker) Stop() {
	close(w.workerQuitChan)
}
