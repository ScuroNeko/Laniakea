package tgapi

import (
	"context"
	"errors"
	"testing"
)

func TestWorkerPoolSubmitAfterStop(t *testing.T) {
	pool := newWorkerPool(1, 1)
	pool.start()
	pool.stop()

	if _, err := pool.submit(context.Background(), func(context.Context) (any, error) {
		return nil, nil
	}); !errors.Is(err, ErrPoolStopped) {
		t.Fatalf("expected ErrPoolStopped, got %v", err)
	}
}

func TestWorkerPoolQueueFull(t *testing.T) {
	pool := newWorkerPool(1, 1)
	pool.start()
	defer pool.stop()

	started := make(chan struct{})
	release := make(chan struct{})

	firstResult, err := pool.submit(context.Background(), func(context.Context) (any, error) {
		close(started)
		<-release
		return "first", nil
	})
	if err != nil {
		t.Fatalf("first submit returned error: %v", err)
	}
	<-started

	secondResult, err := pool.submit(context.Background(), func(context.Context) (any, error) {
		return "second", nil
	})
	if err != nil {
		t.Fatalf("second submit returned error: %v", err)
	}

	if _, err := pool.submit(context.Background(), func(context.Context) (any, error) {
		return "third", nil
	}); !errors.Is(err, ErrPoolQueueFull) {
		t.Fatalf("expected ErrPoolQueueFull, got %v", err)
	}

	close(release)

	first := <-firstResult
	if first.err != nil || first.value != "first" {
		t.Fatalf("unexpected first result: %+v", first)
	}
	second := <-secondResult
	if second.err != nil || second.value != "second" {
		t.Fatalf("unexpected second result: %+v", second)
	}
}
