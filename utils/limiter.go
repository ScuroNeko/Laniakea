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
		globalLimiter: rate.NewLimiter(30, 30),
		chatLimiters:  make(map[int64]*rate.Limiter),
		chatLocks:     make(map[int64]time.Time), // инициализация!
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
	if retryAfter <= 0 {
		return
	}
	rl.chatMu.Lock()
	defer rl.chatMu.Unlock()
	rl.chatLocks[chatID] = time.Now().Add(time.Duration(retryAfter) * time.Second)
}

// GlobalWait блокирует до возможности сделать глобальный запрос.
func (rl *RateLimiter) GlobalWait(ctx context.Context) error {
	// Ждём окончания глобальной блокировки, если она есть
	if err := rl.waitForGlobalUnlock(ctx); err != nil {
		return err
	}
	// Ждём разрешения rate limiter'а
	return rl.globalLimiter.Wait(ctx)
}

// Wait блокирует до возможности сделать запрос для конкретного чата.
func (rl *RateLimiter) Wait(ctx context.Context, chatID int64) error {
	// Ждём окончания блокировки чата
	if err := rl.waitForChatUnlock(ctx, chatID); err != nil {
		return err
	}
	// Затем глобальной блокировки
	if err := rl.waitForGlobalUnlock(ctx); err != nil {
		return err
	}
	// Получаем или создаём лимитер для чата
	limiter := rl.getChatLimiter(chatID)
	return limiter.Wait(ctx)
}

// GlobalAllow неблокирующая проверка глобального запроса.
func (rl *RateLimiter) GlobalAllow() bool {
	rl.globalMu.RLock()
	until := rl.globalLockUntil
	rl.globalMu.RUnlock()

	if !until.IsZero() && time.Now().Before(until) {
		return false
	}
	return rl.globalLimiter.Allow()
}

// Allow неблокирующая проверка запроса для чата.
func (rl *RateLimiter) Allow(chatID int64) bool {
	// Проверяем глобальную блокировку
	rl.globalMu.RLock()
	globalUntil := rl.globalLockUntil
	rl.globalMu.RUnlock()
	if !globalUntil.IsZero() && time.Now().Before(globalUntil) {
		return false
	}

	// Проверяем блокировку чата
	rl.chatMu.Lock()
	chatUntil, ok := rl.chatLocks[chatID]
	rl.chatMu.Unlock()
	if ok && !chatUntil.IsZero() && time.Now().Before(chatUntil) {
		return false
	}

	// Проверяем глобальный лимитер
	if !rl.globalLimiter.Allow() {
		return false
	}

	// Проверяем лимитер чата
	limiter := rl.getChatLimiter(chatID)
	return limiter.Allow()
}

// Внутренние вспомогательные методы

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

func (rl *RateLimiter) waitForChatUnlock(ctx context.Context, chatID int64) error {
	rl.chatMu.Lock()
	until, ok := rl.chatLocks[chatID]
	rl.chatMu.Unlock()

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

func (rl *RateLimiter) getChatLimiter(chatID int64) *rate.Limiter {
	rl.chatMu.Lock()
	defer rl.chatMu.Unlock()
	if lim, ok := rl.chatLimiters[chatID]; ok {
		return lim
	}
	lim := rate.NewLimiter(1, 1) // 1 запрос/сек
	rl.chatLimiters[chatID] = lim
	return lim
}
