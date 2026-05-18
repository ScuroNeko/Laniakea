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

// TestRateLimiterCleanupEvictsIdleChats guards the memory-leak fix: per-chat
// limiter and lastSeen state must be reclaimed by Cleanup once the entry has
// been idle for longer than the threshold, while still-active chats and
// unexpired cooldowns must survive.
func TestRateLimiterCleanupEvictsIdleChats(t *testing.T) {
	rl := NewRateLimiter()

	// Touch chat 1 to make it tracked, then backdate its last-seen marker
	// so it looks idle from Cleanup's perspective.
	if !rl.Allow(1) {
		t.Fatal("expected initial Allow for chat 1 to succeed")
	}
	rl.chatMu.Lock()
	rl.chatLastSeen[1] = time.Now().Add(-time.Hour)
	rl.chatMu.Unlock()

	// Touch chat 2 so it stays "active".
	if !rl.Allow(2) {
		t.Fatal("expected initial Allow for chat 2 to succeed")
	}

	// Expired cooldown should be evicted; future cooldown should survive.
	rl.SetChatLock(10, 1)
	rl.chatMu.Lock()
	rl.chatLocks[10] = time.Now().Add(-time.Second)
	rl.chatLocks[11] = time.Now().Add(time.Hour)
	rl.chatMu.Unlock()

	rl.Cleanup(time.Minute)

	rl.chatMu.RLock()
	defer rl.chatMu.RUnlock()
	if _, ok := rl.chatLimiters[1]; ok {
		t.Fatal("expected idle chat 1 limiter to be evicted")
	}
	if _, ok := rl.chatLastSeen[1]; ok {
		t.Fatal("expected idle chat 1 lastSeen to be evicted")
	}
	if _, ok := rl.chatLimiters[2]; !ok {
		t.Fatal("expected active chat 2 limiter to remain")
	}
	if _, ok := rl.chatLocks[10]; ok {
		t.Fatal("expected expired chat 10 lock to be evicted")
	}
	if _, ok := rl.chatLocks[11]; !ok {
		t.Fatal("expected future chat 11 lock to remain")
	}
}
