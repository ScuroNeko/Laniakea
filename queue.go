package laniakea

import "fmt"

type Queue[T any] struct {
	queue []T
	size  uint64
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
	return q.queue[0]
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.queue) == 0
}

func (q *Queue[T]) IsFull() bool {
	return q.Length() == q.size
}

func (q *Queue[T]) Length() uint64 {
	return uint64(len(q.queue))
}

func (q *Queue[T]) Dequeue() T {
	el := q.queue[0]
	if q.Length() == 1 {
		q.queue = make([]T, 0)
		return el
	}
	q.queue = q.queue[1:]
	return el
}

func (q *Queue[T]) Raw() []T {
	return q.queue
}
