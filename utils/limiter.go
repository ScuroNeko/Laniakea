package utils

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ErrDropOverflow is returned when drop mode rejects a rate-limited request.
var ErrDropOverflow = errors.New("drop overflow limit")

// RateLimiter implements per-chat and global rate limiting with optional blocking.
// It supports two modes:
//   - "drop" mode: immediately reject if limits are exceeded.
//   - "wait" mode: block until capacity is available.
//
// Per-chat limiters are created lazily and accumulate indefinitely. Call Cleanup
// periodically (e.g. from a background runner) to evict idle entries and prevent
// unbounded memory growth in bots that serve many distinct chats.
type RateLimiter struct {
	globalLockUntil time.Time     // global cooldown timestamp (set by API errors)
	globalLimiter   *rate.Limiter // global token bucket (30 req/sec)
	globalMu        sync.RWMutex  // protects globalLockUntil and globalLimiter

	chatLocks    map[int64]time.Time     // per-chat cooldown timestamps
	chatLimiters map[int64]*rate.Limiter // per-chat token buckets (1 req/sec)
	chatLastSeen map[int64]time.Time     // last access timestamp per chat, for Cleanup eviction
	chatMu       sync.RWMutex            // protects chatLocks, chatLimiters, and chatLastSeen
}

// NewRateLimiter creates a new RateLimiter with default limits.
// Global: 30 requests per second, burst 30.
// Per-chat: 1 request per second, burst 1.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		globalLimiter: rate.NewLimiter(30, 30),
		chatLimiters:  make(map[int64]*rate.Limiter),
		chatLocks:     make(map[int64]time.Time),
		chatLastSeen:  make(map[int64]time.Time),
	}
}

// Cleanup removes per-chat limiter state that has not been touched within
// idleThreshold and chat cooldowns whose expiry has already passed.
//
// Safe to call concurrently with Wait/Allow. Intended for periodic invocation
// from a background runner (e.g. once a minute) to bound memory in long-running
// bots that serve many distinct chats.
func (rl *RateLimiter) Cleanup(idleThreshold time.Duration) {
	now := time.Now()
	rl.chatMu.Lock()
	defer rl.chatMu.Unlock()

	for chatID, lastSeen := range rl.chatLastSeen {
		if now.Sub(lastSeen) <= idleThreshold {
			continue
		}
		delete(rl.chatLimiters, chatID)
		delete(rl.chatLastSeen, chatID)
	}
	for chatID, until := range rl.chatLocks {
		if !until.After(now) {
			delete(rl.chatLocks, chatID)
		}
	}
}

// SetGlobalRate overrides global request-per-second limit and burst.
// If rps <= 0, current settings are kept.
func (rl *RateLimiter) SetGlobalRate(rps int) {
	if rps <= 0 {
		return
	}
	rl.globalMu.Lock()
	defer rl.globalMu.Unlock()
	rl.globalLimiter = rate.NewLimiter(rate.Limit(rps), rps)
}

// SetGlobalLock sets a global cooldown period (e.g., after receiving 429 from Telegram).
// If retryAfter <= 0, no lock is applied.
func (rl *RateLimiter) SetGlobalLock(retryAfter int) {
	if retryAfter <= 0 {
		return
	}
	rl.globalMu.Lock()
	defer rl.globalMu.Unlock()
	rl.globalLockUntil = time.Now().Add(time.Duration(retryAfter) * time.Second)
}

// SetChatLock sets a cooldown for a specific chat (e.g., after 429 for that chat).
// If retryAfter <= 0, no lock is applied.
func (rl *RateLimiter) SetChatLock(chatID int64, retryAfter int) {
	if retryAfter <= 0 {
		return
	}
	rl.chatMu.Lock()
	defer rl.chatMu.Unlock()
	rl.chatLocks[chatID] = time.Now().Add(time.Duration(retryAfter) * time.Second)
}

// GlobalWait blocks until a global request can be made.
// Waits for both global cooldown and token bucket availability.
func (rl *RateLimiter) GlobalWait(ctx context.Context) error {
	if err := rl.waitForGlobalUnlock(ctx); err != nil {
		return err
	}
	limiter := rl.getGlobalLimiter()
	if limiter == nil {
		return nil
	}
	return limiter.Wait(ctx)
}

// Wait blocks until a request for the given chat can be made.
// Waits for: chat cooldown → global cooldown → chat token bucket.
// Note: Global limit is checked *before* chat limit to avoid overloading upstream.
func (rl *RateLimiter) Wait(ctx context.Context, chatID int64) error {
	if err := rl.waitForChatUnlock(ctx, chatID); err != nil {
		return err
	}
	if err := rl.waitForGlobalUnlock(ctx); err != nil {
		return err
	}
	limiter := rl.getGlobalLimiter()
	if limiter != nil {
		if err := limiter.Wait(ctx); err != nil {
			return err
		}
	}
	chatLimiter := rl.getChatLimiter(chatID)
	return chatLimiter.Wait(ctx)
}

// Internal helper that returns the global limiter under read lock.
func (rl *RateLimiter) getGlobalLimiter() *rate.Limiter {
	rl.globalMu.RLock()
	defer rl.globalMu.RUnlock()
	return rl.globalLimiter
}

// GlobalAllow checks if a global request can be made without blocking.
// Returns false if either global cooldown is active or token bucket is exhausted.
func (rl *RateLimiter) GlobalAllow() bool {
	rl.globalMu.RLock()
	until := rl.globalLockUntil
	rl.globalMu.RUnlock()

	if !until.IsZero() && time.Now().Before(until) {
		return false
	}
	limiter := rl.getGlobalLimiter()
	if limiter == nil {
		return true
	}
	return limiter.Allow()
}

// Allow checks if a request for the given chat can be made without blocking.
// Returns false if: global cooldown, chat cooldown, global limiter, or chat limiter denies.
// Note: Global limiter is checked before chat limiter — upstream limits take priority.
func (rl *RateLimiter) Allow(chatID int64) bool {
	// Check global cooldown
	rl.globalMu.RLock()
	globalUntil := rl.globalLockUntil
	rl.globalMu.RUnlock()
	if !globalUntil.IsZero() && time.Now().Before(globalUntil) {
		return false
	}

	// Check chat cooldown
	rl.chatMu.RLock()
	chatUntil, ok := rl.chatLocks[chatID]
	rl.chatMu.RUnlock()
	if ok && !chatUntil.IsZero() && time.Now().Before(chatUntil) {
		return false
	}

	// Check global token bucket
	limiter := rl.getGlobalLimiter()
	if limiter != nil && !limiter.Allow() {
		return false
	}

	// Check chat token bucket
	chatLimiter := rl.getChatLimiter(chatID)
	return chatLimiter.Allow()
}

// Check applies rate limiting based on configuration.
// If dropOverflow is true:
//   - Immediately returns ErrDropOverflow if either global or chat limit is exceeded.
//
// Else:
//   - If chatID != 0: waits for chat-specific capacity (including global limit).
//   - If chatID == 0: waits for global capacity only.
//
// chatID == 0 means no specific chat context (e.g., inline query, webhook without chat).
func (rl *RateLimiter) Check(ctx context.Context, dropOverflow bool, chatID int64) error {
	if dropOverflow {
		if chatID != 0 {
			if !rl.Allow(chatID) {

				return ErrDropOverflow
			}
		} else {
			if !rl.GlobalAllow() {
				return ErrDropOverflow
			}
		}
	} else if chatID != 0 {
		if err := rl.Wait(ctx, chatID); err != nil {
			return err
		}
	} else {
		if err := rl.GlobalWait(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Internal helper that waits for the global cooldown to expire.
func (rl *RateLimiter) waitForGlobalUnlock(ctx context.Context) error {
	rl.globalMu.RLock()
	until := rl.globalLockUntil
	rl.globalMu.RUnlock()

	if until.IsZero() || time.Now().After(until) {
		return nil
	}

	select {
	case <-time.After(time.Until(until)):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Internal helper that waits for a chat-specific cooldown to expire.
func (rl *RateLimiter) waitForChatUnlock(ctx context.Context, chatID int64) error {
	rl.chatMu.RLock()
	until, ok := rl.chatLocks[chatID]
	rl.chatMu.RUnlock()

	if !ok || until.IsZero() || time.Now().After(until) {
		return nil
	}

	select {
	case <-time.After(time.Until(until)):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Internal helper that returns or creates a per-chat limiter.
// Updates chatLastSeen so Cleanup can evict idle entries.
func (rl *RateLimiter) getChatLimiter(chatID int64) *rate.Limiter {
	now := time.Now()

	rl.chatMu.RLock()
	lim, ok := rl.chatLimiters[chatID]
	rl.chatMu.RUnlock()
	if ok {
		rl.chatMu.Lock()
		rl.chatLastSeen[chatID] = now
		rl.chatMu.Unlock()
		return lim
	}

	rl.chatMu.Lock()
	defer rl.chatMu.Unlock()
	if lim, ok := rl.chatLimiters[chatID]; ok {
		rl.chatLastSeen[chatID] = now
		return lim
	}
	lim = rate.NewLimiter(1, 1)
	rl.chatLimiters[chatID] = lim
	rl.chatLastSeen[chatID] = now
	return lim
}
