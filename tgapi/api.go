package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
)

type Api struct {
	token  string
	client *http.Client
	Logger *slog.Logger
}

func NewAPI(token string) *Api {
	l := slog.CreateLogger().Level(utils.GetLoggerLevel()).Prefix("API")
	l.AddWriter(l.CreateJsonStdoutWriter())
	client := &http.Client{Timeout: time.Second * 45}
	return &Api{token, client, l}
}
func (api *Api) CloseApi() error {
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
func (r TelegramRequest[R, P]) DoWithContext(ctx context.Context, api *Api) (R, error) {
	var zero R
	data, err := json.Marshal(r.params)
	if err != nil {
		return zero, err
	}
	buf := bytes.NewBuffer(data)

	u := fmt.Sprintf("https://api.telegram.org/bot%s/%s", api.token, r.method)
	if api.Logger != nil {
		api.Logger.Debugln(strings.ReplaceAll(fmt.Sprintf(
			"POST %s %s", u, buf.String(),
		), api.token, "<TOKEN>"))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u, buf)
	if err != nil {
		return zero, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	res, err := api.client.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	reader := io.LimitReader(res.Body, 10<<20)
	data, err = io.ReadAll(reader)
	if err != nil {
		return zero, err
	}

	if api.Logger != nil {
		api.Logger.Debugln(fmt.Sprintf("RES %s %s", r.method, string(data)))
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
func (r TelegramRequest[R, P]) Do(api *Api) (R, error) {
	ctx := context.Background()
	return r.DoWithContext(ctx, api)
}
