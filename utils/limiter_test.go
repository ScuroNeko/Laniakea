package utils

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRateLimiterCheckDropOverflowHonorsGlobalLock(t *testing.T) {
	rl := NewRateLimiter()
	rl.SetGlobalLock(1)

	if err := rl.Check(context.Background(), true, 0); !errors.Is(err, ErrDropOverflow) {
		t.Fatalf("expected ErrDropOverflow, got %v", err)
	}
}

func TestRateLimiterChatLocksAreScopedPerChat(t *testing.T) {
	rl := NewRateLimiter()
	rl.SetChatLock(42, 1)

	if rl.Allow(42) {
		t.Fatal("expected locked chat to be rejected")
	}
	if !rl.Allow(7) {
		t.Fatal("expected unrelated chat to remain allowed")
	}
}

func TestRateLimiterGlobalWaitRespectsContextCancellation(t *testing.T) {
	rl := NewRateLimiter()
	rl.SetGlobalLock(1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := rl.GlobalWait(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
