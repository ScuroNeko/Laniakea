package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
)

type APIOpts struct {
	token         string
	client        *http.Client
	useTestServer bool
	apiUrl        string

	limiter           *utils.RateLimiter
	dropOverflowLimit bool
}

var ErrPoolUnexpected = errors.New("unexpected response from pool")

func NewAPIOpts(token string) *APIOpts {
	return &APIOpts{token: token, client: nil, useTestServer: false, apiUrl: "https://api.telegram.org"}
}
func (opts *APIOpts) SetHTTPClient(client *http.Client) *APIOpts {
	if client != nil {
		opts.client = client
	}
	return opts
}
func (opts *APIOpts) UseTestServer(use bool) *APIOpts {
	opts.useTestServer = use
	return opts
}
func (opts *APIOpts) SetAPIUrl(apiUrl string) *APIOpts {
	if apiUrl != "" {
		opts.apiUrl = apiUrl
	}
	return opts
}
func (opts *APIOpts) SetLimiter(limiter *utils.RateLimiter) *APIOpts {
	opts.limiter = limiter
	return opts
}
func (opts *APIOpts) SetLimiterDrop(b bool) *APIOpts {
	opts.dropOverflowLimit = b
	return opts
}

type API struct {
	token         string
	client        *http.Client
	logger        *slog.Logger
	useTestServer bool
	apiUrl        string

	pool              *WorkerPool
	Limiter           *utils.RateLimiter
	dropOverflowLimit bool
}

func NewAPI(opts *APIOpts) *API {
	l := slog.CreateLogger().Level(utils.GetLoggerLevel()).Prefix("API")
	l.AddWriter(l.CreateJsonStdoutWriter())
	client := opts.client
	if client == nil {
		client = &http.Client{Timeout: time.Second * 45}
	}
	pool := NewWorkerPool(16, 256)
	pool.Start(context.Background())
	return &API{
		opts.token, client, l,
		opts.useTestServer, opts.apiUrl,
		pool, opts.limiter, opts.dropOverflowLimit,
	}
}
func (api *API) CloseApi() error {
	api.pool.Stop()
	return api.logger.Close()
}
func (api *API) GetLogger() *slog.Logger { return api.logger }

type ResponseParameters struct {
	MigrateToChatID *int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      *int   `json:"retry_after,omitempty"`
}
type ApiResponse[R any] struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	Result      R      `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`

	Parameters *ResponseParameters `json:"parameters,omitempty"`
}
type TelegramRequest[R, P any] struct {
	method string
	params P
}

func NewRequest[R, P any](method string, params P) TelegramRequest[R, P] {
	return TelegramRequest[R, P]{method: method, params: params}
}
func (r TelegramRequest[R, P]) doRequest(ctx context.Context, api *API) (R, error) {
	var zero R
	if api.Limiter != nil {
		if api.dropOverflowLimit {
			if !api.Limiter.GlobalAllow() {
				return zero, errors.New("rate limited")
			}
		} else {
			if err := api.Limiter.GlobalWait(ctx); err != nil {
				return zero, err
			}
		}
	}

	data, err := json.Marshal(r.params)
	if err != nil {
		return zero, err
	}
	buf := bytes.NewBuffer(data)

	methodPrefix := ""
	if api.useTestServer {
		methodPrefix = "/test"
	}
	url := fmt.Sprintf("%s/bot%s%s/%s", api.apiUrl, api.token, methodPrefix, r.method)
	req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		return zero, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	api.logger.Debugln("REQ", api.apiUrl, r.method, buf.String())
	res, err := api.client.Do(req)
	if err != nil {
		return zero, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	data, err = readBody(res.Body)
	if err != nil {
		return zero, err
	}
	api.logger.Debugln("RES", r.method, string(data))
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusTooManyRequests {
		return zero, fmt.Errorf("unexpected status code: %d, %s", res.StatusCode, string(data))
	}

	responseData, err := parseBody[R](data)
	if errors.Is(err, ErrRateLimit) {
		if responseData.Parameters != nil {
			after := 0
			if responseData.Parameters.RetryAfter != nil {
				after = *responseData.Parameters.RetryAfter
			}
			api.Limiter.SetGlobalLock(after)
			return r.doRequest(ctx, api)
		}
		return zero, ErrRateLimit
	}
	return responseData.Result, err
}
func (r TelegramRequest[R, P]) DoWithContext(ctx context.Context, api *API) (R, error) {
	var zero R
	result, err := api.pool.Submit(ctx, func(ctx context.Context) (any, error) {
		return r.doRequest(ctx, api)
	})
	if err != nil {
		return zero, err
	}

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case res := <-result:
		if res.Err != nil {
			return zero, res.Err
		}
		if val, ok := res.Value.(R); ok {
			return val, nil
		}
		return zero, ErrPoolUnexpected
	}
}
func (r TelegramRequest[R, P]) Do(api *API) (R, error) {
	return r.DoWithContext(context.Background(), api)
}

func readBody(body io.ReadCloser) ([]byte, error) {
	reader := io.LimitReader(body, 10<<20)
	return io.ReadAll(reader)
}
func parseBody[R any](data []byte) (ApiResponse[R], error) {
	var resp ApiResponse[R]
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return resp, err
	}
	if !resp.Ok {
		if resp.ErrorCode == 429 {
			return resp, ErrRateLimit
		}
		return resp, fmt.Errorf("[%d] %s", resp.ErrorCode, resp.Description)
	}
	return resp, nil
}
