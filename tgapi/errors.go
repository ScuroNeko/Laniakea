package tgapi

import "errors"

var ErrRateLimit = errors.New("rate limit exceeded")
var ErrPoolUnexpected = errors.New("unexpected response from pool")
var ErrPoolQueueFull = errors.New("worker pool queue full")
