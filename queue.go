package laniakea

import (
	"fmt"
	"sync"
)

type Queue[T any] struct {
	size  uint64
	mu    sync.RWMutex
	queue []T
}

func CreateQueue[T any](size uint64) *Queue[T] {
	return &Queue[T]{
		queue: make([]T, 0),
		size:  size,
	}
}

func (q *Queue[T]) Enqueue(el T) error {
	if q.IsFull() {
		return fmt.Errorf("queue full")
	}
	q.queue = append(q.queue, el)
	return nil
}

func (q *Queue[T]) Peak() T {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.queue[0]
}

func (q *Queue[T]) IsEmpty() bool {
	return q.Length() == 0
}

func (q *Queue[T]) IsFull() bool {
	return q.Length() == q.size
}

func (q *Queue[T]) Length() uint64 {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return uint64(len(q.queue))
}

func (q *Queue[T]) Dequeue() T {
	q.mu.RLock()
	el := q.queue[0]
	q.mu.RUnlock()

	if q.Length() == 1 {
		q.mu.Lock()
		q.queue = make([]T, 0)
		q.mu.Unlock()
		return el
	}

	q.mu.Lock()
	q.queue = q.queue[1:]
	q.mu.Unlock()
	return el
}

func (q *Queue[T]) Raw() []T {
	return q.queue
}
