package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/utils"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

// APIOpts holds configuration options for initializing the Telegram API client.
// Use the provided setter methods to build options — do not construct directly.
type APIOpts struct {
	token         string
	client        *http.Client
	useTestServer bool
	apiURL        string

	logFormat    utils.LogFormat
	logFormatter *sneklog.Formatter

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
		apiURL:        "https://api.telegram.org",
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

// SetAPIURL overrides the default Telegram API URL.
// Useful for self-hosted bots or proxies.
func (opts *APIOpts) SetAPIURL(apiURL string) *APIOpts {
	if apiURL != "" {
		opts.apiURL = apiURL
	}
	return opts
}

// SetLogFormat sets the output format used by API-managed loggers.
func (opts *APIOpts) SetLogFormat(format utils.LogFormat) *APIOpts {
	opts.logFormat = format
	return opts
}

// SetLogFormatter sets the formatter used by API-managed logger writers.
func (opts *APIOpts) SetLogFormatter(formatter *sneklog.Formatter) *APIOpts {
	opts.logFormatter = formatter
	return opts
}

// SetLimiter sets a rate limiter to enforce Telegram's API limits.
// Recommended: use utils.NewRateLimiter() for correct per-chat and global throttling.
func (opts *APIOpts) SetLimiter(limiter *utils.RateLimiter) *APIOpts {
	opts.limiter = limiter
	return opts
}

// SetDropRateLimitOverflow enables "drop mode" for rate limiting.
// If true, requests exceeding limits return ErrDropOverflow immediately.
// If false, requests block until capacity is available.
func (opts *APIOpts) SetDropRateLimitOverflow(b bool) *APIOpts {
	opts.dropOverflowLimit = b
	return opts
}

// API is the main Telegram Bot API client for JSON requests.
//
// Use API methods when sending JSON payloads (for example with file_id, URL, or other
// non-multipart fields). For multipart file uploads, use Uploader.
//
// It manages HTTP requests, rate limiting, retries, and connection pooling.
type API struct {
	token         string
	client        *http.Client
	logger        *sneklog.Logger
	useTestServer bool
	apiURL        string

	logFormat    utils.LogFormat
	logFormatter *sneklog.Formatter

	pool              *workerPool
	Limiter           *utils.RateLimiter
	dropOverflowLimit bool
}

// NewAPI creates a new API client from options.
// Always call Close() when done to release resources.
func NewAPI(opts *APIOpts) *API {
	if opts == nil {
		return nil
	}
	logger := utils.CreateLogger(
		"API", utils.GetLoggerLevel(),
		opts.logFormat, opts.logFormatter,
	)

	client := opts.client
	if client == nil {
		client = &http.Client{Timeout: time.Second * 45}
	}

	pool := newWorkerPool(16, 256)
	pool.start()

	return &API{
		token:         opts.token,
		client:        client,
		logger:        logger,
		useTestServer: opts.useTestServer,
		apiURL:        opts.apiURL,

		logFormat:    opts.logFormat,
		logFormatter: opts.logFormatter,

		pool:              pool,
		Limiter:           opts.limiter,
		dropOverflowLimit: opts.dropOverflowLimit,
	}
}

// Close shuts down the internal worker pool and closes the logger.
// Must be called to avoid resource leaks.
// See https://core.telegram.org/bots/api
func (api *API) Close() error {
	api.pool.stop()
	if api.client != nil {
		api.client.CloseIdleConnections()
	}
	return api.logger.Close()
}

// GetLogger returns the internal logger for custom logging.
// See https://core.telegram.org/bots/api
func (api *API) GetLogger() *sneklog.Logger {
	return api.logger
}

// ResponseParameters contains Telegram API response metadata (e.g., retry_after, migrate_to_chat_id).
type ResponseParameters struct {
	MigrateToChatID *int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      *int   `json:"retry_after,omitempty"`
}

// TelegramResponse is the standard Telegram Bot API response structure.
// Generic over Result type R.
type TelegramResponse[R any] struct {
	Ok          bool                `json:"ok"`
	Description string              `json:"description,omitempty"`
	Result      R                   `json:"result,omitempty"`
	ErrorCode   int                 `json:"error_code,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

// TelegramRequest is a low-level Telegram API request wrapper.
//
// Prefer method-specific helpers such as SendMessage or GetUpdates. TelegramRequest
// bypasses method-specific parameter types and convenience helpers, so callers are
// responsible for using the correct method name and compatible request and response types.
// In that sense it is an unsafe escape hatch compared with the typed API surface.
type TelegramRequest[R, P any] struct {
	method string
	params P
	chatID int64
}

// NewRequest creates a low-level TelegramRequest with no associated chat ID.
func NewRequest[R, P any](method string, params P) TelegramRequest[R, P] {
	return TelegramRequest[R, P]{method, params, 0}
}

// NewRequestWithChatID creates a low-level TelegramRequest with an associated chat ID.
// The chat ID is used for per-chat rate limiting.
func NewRequestWithChatID[R, P any](method string, params P, chatID int64) TelegramRequest[R, P] {
	return TelegramRequest[R, P]{method, params, chatID}
}

func (r TelegramRequest[R, P]) doRequest(ctx context.Context, api *API) (R, error) {
	var zero R
	reqData, err := json.Marshal(r.params)
	if err != nil {
		return zero, fmt.Errorf("failed to marshal request: %w", err)
	}

	methodPrefix := ""
	if api.useTestServer {
		methodPrefix = "/test"
	}
	url := fmt.Sprintf("%s/bot%s%s/%s", api.apiURL, api.token, methodPrefix, r.method)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return zero, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	for {
		// Apply rate limiting before making the request
		if api.Limiter != nil {
			if err := api.Limiter.Check(ctx, api.dropOverflowLimit, r.chatID); err != nil {
				return zero, err
			}
		}
		buf := bytes.NewBuffer(reqData)
		req.Body = io.NopCloser(buf)
		req.ContentLength = int64(len(reqData))

		api.logger.Debugln("REQ", url, string(reqData))
		resp, err := api.client.Do(req)
		if err != nil {
			return zero, fmt.Errorf("HTTP request failed: %w", err)
		}

		respData, err := readBody(resp.Body)
		_ = resp.Body.Close() // ensure body is closed
		if err != nil {
			return zero, fmt.Errorf("failed to read response body: %w", err)
		}

		api.logger.Debugln("RES", r.method, string(respData))

		response, err := parseBody[R](respData)
		if err != nil {
			return zero, fmt.Errorf("failed to parse response: %w", err)
		}

		if !response.Ok {
			// Handle rate limiting (429)
			if response.ErrorCode == 429 && response.Parameters != nil && response.Parameters.RetryAfter != nil {
				after := *response.Parameters.RetryAfter
				api.logger.Warnf("Rate limited by Telegram, retry after %d seconds (chat: %d)", after, r.chatID)

				// Apply cooldown to global or chat-specific limiter
				if api.Limiter != nil {
					if r.chatID > 0 {
						api.Limiter.SetChatLock(r.chatID, after)
					} else {
						api.Limiter.SetGlobalLock(after)
					}
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

// Internal helper that reads and caps a Telegram response body.
func readBody(body io.ReadCloser) ([]byte, error) {
	reader := io.LimitReader(body, 10<<20) // 10 MB
	return io.ReadAll(reader)
}

// Internal helper that parses a typed Telegram API response body.
func parseBody[R any](data []byte) (TelegramResponse[R], error) {
	var resp TelegramResponse[R]
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return resp, nil
}
