package laniakea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.nix13.pw/scuroneko/slog"
)

type Api struct {
	token  string
	logger *slog.Logger
}

func NewAPI(token string) *Api {
	l := slog.CreateLogger().Level(GetLoggerLevel()).Prefix("API")
	l.AddWriter(l.CreateJsonStdoutWriter())
	return &Api{token, l}
}
func (api *Api) CloseApi() error {
	return api.logger.Close()
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
func (r TelegramRequest[R, P]) Do(api *Api) (*R, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(r.params)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("https://api.telegram.org/bot%s/%s", api.token, r.method)
	if api.logger != nil {
		api.logger.Debugln(strings.ReplaceAll(fmt.Sprintf(
			"POST %s %s", u, buf.String(),
		), api.token, "<TOKEN>"))
	}

	res, err := http.Post(u, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if api.logger != nil {
		api.logger.Debugln(fmt.Sprintf("RES %s %s", r.method, string(data)))
	}

	response := new(ApiResponse[R])
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
	}
	return &response.Result, nil
}

func (b *Bot) GetFileByLink(link string) ([]byte, error) {
	u := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.token, link)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
