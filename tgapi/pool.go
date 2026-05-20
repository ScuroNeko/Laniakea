package tgapi

import (
	"context"
	"sync"
)

type workerPool struct {
	taskCh    chan requestEnvelope
	queueSize int
	workers   int
	wg        sync.WaitGroup
	quit      chan struct{}
	stopOnce  sync.Once
	started   bool
	stopped   bool
	startedMu sync.Mutex
}

type requestEnvelope struct {
	ctx      context.Context
	doFunc   func(context.Context) (any, error)
	resultCh chan requestResult
}

type requestResult struct {
	value any
	err   error
}

func newWorkerPool(workers int, queueSize int) *workerPool {
	if workers <= 0 {
		workers = 1
	}
	if queueSize <= 0 {
		queueSize = 100
	}

	return &workerPool{
		taskCh:    make(chan requestEnvelope, queueSize),
		queueSize: queueSize,
		workers:   workers,
		quit:      make(chan struct{}),
	}
}

func (p *workerPool) start() {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()
	if p.started {
		return
	}
	p.started = true

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *workerPool) stop() {
	p.stopOnce.Do(func() {
		p.startedMu.Lock()
		p.stopped = true
		p.started = false
		close(p.quit)
		p.startedMu.Unlock()

		p.wg.Wait()
	})
}

func (p *workerPool) submit(ctx context.Context, do func(context.Context) (any, error)) (<-chan requestResult, error) {
	p.startedMu.Lock()
	if p.stopped || !p.started {
		p.startedMu.Unlock()
		return nil, ErrPoolStopped
	}

	if len(p.taskCh) >= p.queueSize {
		p.startedMu.Unlock()
		return nil, ErrPoolQueueFull
	}

	resultCh := make(chan requestResult, 1)

	envelope := requestEnvelope{
		ctx:      ctx,
		doFunc:   do,
		resultCh: resultCh,
	}

	select {
	case <-ctx.Done():
		p.startedMu.Unlock()
		return nil, ctx.Err()
	case p.taskCh <- envelope:
		p.startedMu.Unlock()
		return resultCh, nil
	default:
		p.startedMu.Unlock()
		return nil, ErrPoolQueueFull
	}
}

func (p *workerPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.quit:
			// Drain queued work after stop. No new tasks are accepted.
			for {
				select {
				case envelope := <-p.taskCh:
					p.executeEnvelope(envelope)
				default:
					return
				}
			}

		case envelope := <-p.taskCh:
			p.executeEnvelope(envelope)
		}
	}
}

func (p *workerPool) executeEnvelope(envelope requestEnvelope) {
	value, err := envelope.doFunc(envelope.ctx)
	envelope.resultCh <- requestResult{
		value: value,
		err:   err,
	}
	close(envelope.resultCh)
}
