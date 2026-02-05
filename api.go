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
	return &Api{
		token:  token,
		logger: l,
	}
}
func (api *Api) CloseApi() {
	api.logger.Close()
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
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r.params)
	if err != nil {
		return nil, err
	}

	if api.logger != nil {
		api.logger.Debugln(strings.ReplaceAll(fmt.Sprintf(
			"POST https://api.telegram.org/bot%s/%s %s",
			"<TOKEN>", r.method, buf.String(),
		), "\n", ""))
	}

	req, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/%s", api.token, r.method), "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	data, err := io.ReadAll(req.Body)
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
	c := http.DefaultClient
	u := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.token, link)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
