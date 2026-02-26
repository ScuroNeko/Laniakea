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
	"golang.org/x/time/rate"
)

type APIOpts struct {
	token         string
	client        *http.Client
	useTestServer bool
	apiUrl        string

	limiter           *rate.Limiter
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
func (opts *APIOpts) SetLimiter(limiter *rate.Limiter) *APIOpts {
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
	limiter           *rate.Limiter
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

type ApiResponse[R any] struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	Result      R      `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
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
	if api.limiter != nil {
		if api.dropOverflowLimit {
			if !api.limiter.Allow() {
				return zero, errors.New("rate limited")
			}
		} else {
			if err := api.limiter.Wait(ctx); err != nil {
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
	if res.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("unexpected status code: %d, %s", res.StatusCode, string(data))
	}
	return parseBody[R](data)
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
func parseBody[R any](data []byte) (R, error) {
	var zero R
	var resp ApiResponse[R]
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return zero, err
	}
	if !resp.Ok {
		return zero, fmt.Errorf("[%d] %s", resp.ErrorCode, resp.Description)
	}
	return resp.Result, nil
}
