package tgapi

import "errors"

// ErrRateLimit reports that a request exceeded the configured rate limiter.
var ErrRateLimit = errors.New("rate limit exceeded")

// ErrPoolUnexpected reports an unexpected result type returned from the worker pool.
var ErrPoolUnexpected = errors.New("unexpected response from pool")

// ErrPoolQueueFull reports that the internal request queue is full.
var ErrPoolQueueFull = errors.New("worker pool queue full")

// ErrPoolStopped reports that a request was submitted after the worker pool stopped.
var ErrPoolStopped = errors.New("worker pool stopped")
