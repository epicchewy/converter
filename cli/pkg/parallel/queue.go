package parallel

import "context"

type Queue interface {
	Enqueue(ctx context.Context, i interface{}) error
	Dequeue(ctx context.Context) (interface{}, error)
	Size() int
	IsEmpty() bool
	Stop()
}

type ChanQueue struct {
	queue chan interface{}
	size  int
	stop  chan struct{}
}

func NewChanQueue(ctx context.Context, maxSize int) *ChanQueue {
	return &ChanQueue{
		queue: make(chan interface{}, maxSize),
		size:  maxSize,
		stop:  make(chan struct{}),
	}
}

func (q *ChanQueue) Enqueue(ctx context.Context, i interface{}) error {
	select {
	case q.queue <- i:
		return nil
	case <-q.stop:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *ChanQueue) Dequeue(ctx context.Context) (interface{}, error) {
	select {
	case i := <-q.queue:
		return i, nil
	case <-q.stop:
		return nil, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (q *ChanQueue) IsEmpty() bool {
	return len(q.queue) == 0
}

func (q *ChanQueue) Size() int {
	return len(q.queue)
}

func (q *ChanQueue) Stop() {
	close(q.stop)
}
