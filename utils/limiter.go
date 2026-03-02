package utils

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	globalLockUntil time.Time
	globalLimiter   *rate.Limiter
	globalMu        sync.RWMutex

	chatLocks    map[int64]time.Time
	chatLimiters map[int64]*rate.Limiter
	chatMu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		// 30 запросов в секунду (burst=30)
		globalLimiter: rate.NewLimiter(rate.Limit(30), 30),
		chatLimiters:  make(map[int64]*rate.Limiter),
	}
}

func (rl *RateLimiter) SetGlobalLock(retryAfter int) {
	if retryAfter <= 0 {
		return
	}
	rl.globalMu.Lock()
	defer rl.globalMu.Unlock()
	rl.globalLockUntil = time.Now().Add(time.Duration(retryAfter) * time.Second)
}
func (rl *RateLimiter) SetChatLock(chatID int64, retryAfter int) {
	rl.chatMu.Lock()
	defer rl.chatMu.Unlock()
	rl.chatLocks[chatID] = time.Now().Add(time.Duration(retryAfter) * time.Second)
}

func (rl *RateLimiter) GlobalWait(ctx context.Context) error {
	rl.globalMu.RLock()
	until := rl.globalLockUntil
	rl.globalMu.RUnlock()

	if !until.IsZero() {
		if time.Now().Before(until) {
			// Ждём до окончания блокировки или отмены контекста
			select {
			case <-time.After(time.Until(until)):
				// блокировка снята
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	// Теперь ждём разрешения rate limiter'а
	return rl.globalLimiter.Wait(ctx)
}
func (rl *RateLimiter) Wait(ctx context.Context, chatID int64) error {
	rl.chatMu.Lock()
	until, ok := rl.chatLocks[chatID]
	rl.chatMu.Unlock()
	if ok && !until.IsZero() {
		if time.Now().Before(until) {
			select {
			case <-time.After(time.Until(until)):
				// блокировка снята
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	if err := rl.GlobalWait(ctx); err != nil {
		return err
	}
	rl.chatMu.Lock()
	chatLimiter, ok := rl.chatLimiters[chatID]
	if !ok {
		chatLimiter = rate.NewLimiter(rate.Limit(1), 1)
		rl.chatLimiters[chatID] = chatLimiter
	}
	rl.chatMu.Unlock()
	return chatLimiter.Wait(ctx)
}

func (rl *RateLimiter) GlobalAllow() bool {
	rl.globalMu.RLock()
	until := rl.globalLockUntil
	rl.globalMu.RUnlock()

	if !until.IsZero() {
		if time.Now().Before(until) {
			// Ждём до окончания блокировки или отмены контекста
			select {
			case <-time.After(time.Until(until)):
				rl.globalLimiter.Allow()
			}
		}
	}
	return rl.globalLimiter.Allow()
}
func (rl *RateLimiter) Allow(chatID int64) bool {
	rl.chatMu.Lock()
	until, ok := rl.chatLocks[chatID]
	rl.chatMu.Unlock()
	if ok && !until.IsZero() {
		if time.Now().Before(until) {
			select {
			case <-time.After(time.Until(until)):
			}
		}
	}

	if !rl.globalLimiter.Allow() {
		return false
	}
	rl.chatMu.Lock()
	chatLimiter, ok := rl.chatLimiters[chatID]
	if !ok {
		chatLimiter = rate.NewLimiter(rate.Limit(1), 1)
		rl.chatLimiters[chatID] = chatLimiter
	}
	rl.chatMu.Unlock()
	return chatLimiter.Allow()
}
