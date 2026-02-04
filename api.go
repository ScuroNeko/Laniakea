package laniakea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
func (r TelegramRequest[R, P]) Do(bot *Bot) (*R, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r.params)
	if err != nil {
		return nil, err
	}

	if bot.requestLogger != nil {
		bot.requestLogger.Debugln(strings.ReplaceAll(fmt.Sprintf(
			"POST https://api.telegram.org/bot%s/%s %s",
			"<TOKEN>", r.method, buf.String(),
		), "\n", ""))
	}

	req, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.token, r.method), "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	if bot.requestLogger != nil {
		bot.requestLogger.Debugln(fmt.Sprintf("RES %s %s", r.method, string(data)))
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
