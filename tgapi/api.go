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

type API struct {
	token  string
	client *http.Client
	Logger *slog.Logger
}

func NewAPI(token string) *API {
	l := slog.CreateLogger().Level(utils.GetLoggerLevel()).Prefix("API")
	l.AddWriter(l.CreateJsonStdoutWriter())
	client := &http.Client{Timeout: time.Second * 45}
	return &API{token, client, l}
}
func (api *API) CloseApi() error {
	return api.Logger.Close()
}

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
func (r TelegramRequest[R, P]) DoWithContext(ctx context.Context, api *API) (R, error) {
	var zero R
	data, err := json.Marshal(r.params)
	if err != nil {
		return zero, err
	}
	buf := bytes.NewBuffer(data)

	u := fmt.Sprintf("https://api.telegram.org/bot%s/%s", api.token, r.method)
	req, err := http.NewRequestWithContext(ctx, "POST", u, buf)
	if err != nil {
		return zero, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	api.Logger.Debugln("REQ", r.method, buf.String())
	res, err := api.client.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	reader := io.LimitReader(res.Body, 10<<20)
	data, err = io.ReadAll(reader)
	if err != nil {
		return zero, err
	}
	api.Logger.Debugln("RES", r.method, string(data))
	if res.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("unexpected status code: %d, %s", res.StatusCode, string(data))
	}

	var resp ApiResponse[R]
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return zero, err
	}
	if !resp.Ok {
		return zero, fmt.Errorf("[%d] %s", resp.ErrorCode, resp.Description)
	}
	return resp.Result, nil

}
func (r TelegramRequest[R, P]) Do(api *API) (R, error) {
	ctx := context.Background()
	return r.DoWithContext(ctx, api)
}
