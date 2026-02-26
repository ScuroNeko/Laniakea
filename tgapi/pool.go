package tgapi

import (
	"context"
	"errors"
	"sync"
)

var ErrPoolQueueFull = errors.New("worker pool queue full")

type RequestEnvelope struct {
	DoFunc   func(context.Context) (any, error) // функция, которая выполнит запрос и вернет any
	ResultCh chan RequestResult                 // канал для результата
}
type RequestResult struct {
	Value any
	Err   error
}

// WorkerPool управляет воркерами и очередью
type WorkerPool struct {
	taskCh    chan RequestEnvelope
	queueSize int
	workers   int
	wg        sync.WaitGroup
	quit      chan struct{}
	started   bool
	startedMu sync.Mutex
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		taskCh:    make(chan RequestEnvelope, queueSize),
		queueSize: queueSize,
		workers:   workers,
		quit:      make(chan struct{}),
	}
}

// Start запускает воркеров
func (p *WorkerPool) Start(ctx context.Context) {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()
	if p.started {
		return
	}
	p.started = true
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

// Stop останавливает пул (ждет завершения текущих задач)
func (p *WorkerPool) Stop() {
	close(p.quit)
	p.wg.Wait()
}

// Submit отправляет задачу в очередь и возвращает канал для результата
func (p *WorkerPool) Submit(ctx context.Context, do func(context.Context) (any, error)) (<-chan RequestResult, error) {
	if len(p.taskCh) >= p.queueSize {
		return nil, ErrPoolQueueFull
	}

	resultCh := make(chan RequestResult, 1) // буфер 1, чтобы не блокировать воркера
	envelope := RequestEnvelope{do, resultCh}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case p.taskCh <- envelope:
		return resultCh, nil
	default:
		return nil, ErrPoolQueueFull
	}
}

// worker выполняет задачи
func (p *WorkerPool) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case <-p.quit:
			return
		case envelope := <-p.taskCh:
			// Выполняем задачу с переданным контекстом (или можно использовать свой)
			val, err := envelope.DoFunc(ctx)
			envelope.ResultCh <- RequestResult{Value: val, Err: err}
			close(envelope.ResultCh)
		}
	}
}
