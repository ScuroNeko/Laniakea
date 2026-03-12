package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
)

// APIOpts holds configuration options for initializing the Telegram API client.
// Use the provided setter methods to build options — do not construct directly.
type APIOpts struct {
	token         string
	client        *http.Client
	useTestServer bool
	apiUrl        string

	limiter           *utils.RateLimiter
	dropOverflowLimit bool
}

// NewAPIOpts creates a new APIOpts with default values.
// Use setter methods to customize behavior.
func NewAPIOpts(token string) *APIOpts {
	return &APIOpts{
		token:         token,
		client:        nil,
		useTestServer: false,
		apiUrl:        "https://api.telegram.org",
	}
}

// SetHTTPClient sets a custom HTTP client. Use this for timeouts, proxies, or custom transport.
// If not set, a default client with 45s timeout is used.
func (opts *APIOpts) SetHTTPClient(client *http.Client) *APIOpts {
	if client != nil {
		opts.client = client
	}
	return opts
}

// UseTestServer enables use of Telegram's test server (https://api.test.telegram.org).
// Only for development/testing.
func (opts *APIOpts) UseTestServer(use bool) *APIOpts {
	opts.useTestServer = use
	return opts
}

// SetAPIUrl overrides the default Telegram API URL.
// Useful for self-hosted bots or proxies.
func (opts *APIOpts) SetAPIUrl(apiUrl string) *APIOpts {
	if apiUrl != "" {
		opts.apiUrl = apiUrl
	}
	return opts
}

// SetLimiter sets a rate limiter to enforce Telegram's API limits.
// Recommended: use utils.NewRateLimiter() for correct per-chat and global throttling.
func (opts *APIOpts) SetLimiter(limiter *utils.RateLimiter) *APIOpts {
	opts.limiter = limiter
	return opts
}

// SetLimiterDrop enables "drop mode" for rate limiting.
// If true, requests exceeding limits return ErrDropOverflow immediately.
// If false, requests block until capacity is available.
func (opts *APIOpts) SetLimiterDrop(b bool) *APIOpts {
	opts.dropOverflowLimit = b
	return opts
}

// API is the main Telegram Bot API client.
// It manages HTTP requests, rate limiting, retries, and connection pooling.
type API struct {
	token         string
	client        *http.Client
	logger        *slog.Logger
	useTestServer bool
	apiUrl        string

	pool              *workerPool
	Limiter           *utils.RateLimiter
	dropOverflowLimit bool
}

// NewAPI creates a new API client from options.
// Always call CloseApi() when done to release resources.
func NewAPI(opts *APIOpts) *API {
	l := slog.CreateLogger().Level(utils.GetLoggerLevel()).Prefix("API")
	l.AddWriter(l.CreateJsonStdoutWriter())

	client := opts.client
	if client == nil {
		client = &http.Client{Timeout: time.Second * 45}
	}

	pool := newWorkerPool(16, 256)
	pool.start(context.Background())

	return &API{
		token:             opts.token,
		client:            client,
		logger:            l,
		useTestServer:     opts.useTestServer,
		apiUrl:            opts.apiUrl,
		pool:              pool,
		Limiter:           opts.limiter,
		dropOverflowLimit: opts.dropOverflowLimit,
	}
}

// CloseApi shuts down the internal worker pool and closes the logger.
// Must be called to avoid resource leaks.
func (api *API) CloseApi() error {
	api.pool.stop()
	return api.logger.Close()
}

// GetLogger returns the internal logger for custom logging.
func (api *API) GetLogger() *slog.Logger {
	return api.logger
}

// ResponseParameters contains Telegram API response metadata (e.g., retry_after, migrate_to_chat_id).
type ResponseParameters struct {
	MigrateToChatID *int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      *int   `json:"retry_after,omitempty"`
}

// ApiResponse is the standard Telegram Bot API response structure.
// Generic over Result type R.
type ApiResponse[R any] struct {
	Ok          bool                `json:"ok"`
	Description string              `json:"description,omitempty"`
	Result      R                   `json:"result,omitempty"`
	ErrorCode   int                 `json:"error_code,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

// TelegramRequest is an internal helper struct.
// DO NOT USE NewRequest or NewRequestWithChatID — they are unsafe and discouraged.
// Instead, use explicit methods like SendMessage, GetUpdates, etc.
//
// Why? Because using generics with arbitrary types P and R leads to:
//   - No compile-time validation of parameters
//   - No IDE autocompletion
//   - Runtime panics on malformed JSON
//   - Hard-to-debug errors
//
// Recommended: Define specific methods for each Telegram method (see below).
type TelegramRequest[R, P any] struct {
	method string
	params P
	chatId int64
}

// NewRequest and NewRequestWithChatID are DEPRECATED.
// They encourage unsafe, untyped usage and bypass Go's type safety.
// Instead, define explicit, type-safe methods for each Telegram API endpoint.
//
// Example:
//
//	func (api *API) SendMessage(ctx context.Context, chatID int64, text string) (Message, error) { ... }
//
// This provides:
//
//	✅ Compile-time validation
//	✅ IDE autocompletion
//	✅ Clear API surface
//	✅ Better error messages
//
// DO NOT use these constructors in production code.
// This can be used ONLY for testing or if you NEED method, that wasn't added as function.
func NewRequest[R, P any](method string, params P) TelegramRequest[R, P] {
	return TelegramRequest[R, P]{method, params, 0}
}

func NewRequestWithChatID[R, P any](method string, params P, chatId int64) TelegramRequest[R, P] {
	return TelegramRequest[R, P]{method, params, chatId}
}

// doRequest performs a single HTTP request to Telegram API.
// Handles rate limiting, retries on 429, and parses responses.
// Must be called within a worker pool context if using DoWithContext.
func (r TelegramRequest[R, P]) doRequest(ctx context.Context, api *API) (R, error) {
	var zero R

	data, err := json.Marshal(r.params)
	if err != nil {
		return zero, fmt.Errorf("failed to marshal request: %w", err)
	}
	buf := bytes.NewBuffer(data)

	methodPrefix := ""
	if api.useTestServer {
		methodPrefix = "/test"
	}
	url := fmt.Sprintf("%s/bot%s%s/%s", api.apiUrl, api.token, methodPrefix, r.method)

	req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		return zero, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))
	req.Header.Set("Accept-Encoding", "gzip")
	req.ContentLength = int64(len(data))

	for {
		// Apply rate limiting before making the request
		if api.Limiter != nil {
			if err := api.Limiter.Check(ctx, api.dropOverflowLimit, r.chatId); err != nil {
				return zero, err
			}
		}

		api.logger.Debugln("REQ", url, string(data))
		resp, err := api.client.Do(req)
		if err != nil {
			return zero, fmt.Errorf("HTTP request failed: %w", err)
		}

		data, err = readBody(resp.Body)
		_ = resp.Body.Close() // ensure body is closed
		if err != nil {
			return zero, fmt.Errorf("failed to read response body: %w", err)
		}

		api.logger.Debugln("RES", r.method, string(data))

		response, err := parseBody[R](data)
		if err != nil {
			return zero, fmt.Errorf("failed to parse response: %w", err)
		}

		if !response.Ok {
			// Handle rate limiting (429)
			if response.ErrorCode == 429 && response.Parameters != nil && response.Parameters.RetryAfter != nil {
				after := *response.Parameters.RetryAfter
				api.logger.Warnf("Rate limited by Telegram, retry after %d seconds (chat: %d)", after, r.chatId)

				// Apply cooldown to global or chat-specific limiter
				if r.chatId > 0 {
					api.Limiter.SetChatLock(r.chatId, after)
				} else {
					api.Limiter.SetGlobalLock(after)
				}

				// Wait and retry
				select {
				case <-ctx.Done():
					return zero, ctx.Err()
				case <-time.After(time.Duration(after) * time.Second):
					continue // retry request
				}
			}

			// Other API errors
			return zero, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
		}

		return response.Result, nil
	}
}

// DoWithContext executes the request asynchronously via the worker pool.
// Returns result or error via channel. Respects context cancellation.
func (r TelegramRequest[R, P]) DoWithContext(ctx context.Context, api *API) (R, error) {
	var zero R

	resultChan, err := api.pool.submit(ctx, func(ctx context.Context) (any, error) {
		return r.doRequest(ctx, api)
	})
	if err != nil {
		return zero, err
	}

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case res := <-resultChan:
		if res.err != nil {
			return zero, res.err
		}
		if val, ok := res.value.(R); ok {
			return val, nil
		}
		return zero, ErrPoolUnexpected
	}
}

// Do executes the request synchronously with a background context.
// Use only for simple, non-critical calls.
func (r TelegramRequest[R, P]) Do(api *API) (R, error) {
	return r.DoWithContext(context.Background(), api)
}

// readBody reads and limits response body to prevent memory exhaustion.
// Telegram responses are typically small (<1MB), but we cap at 10MB.
func readBody(body io.ReadCloser) ([]byte, error) {
	reader := io.LimitReader(body, 10<<20) // 10 MB
	return io.ReadAll(reader)
}

// parseBody unmarshals Telegram API response and returns structured result.
// Returns ErrRateLimit internally if error_code == 429 — caller must handle via response.Ok check.
func parseBody[R any](data []byte) (ApiResponse[R], error) {
	var resp ApiResponse[R]
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !resp.Ok {
		if resp.ErrorCode == 429 {
			return resp, ErrRateLimit // internal use only
		}
		return resp, fmt.Errorf("[%d] %s", resp.ErrorCode, resp.Description)
	}

	return resp, nil
}
