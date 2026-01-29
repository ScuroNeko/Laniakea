package laniakea

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ApiResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	Result      any    `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
}
type ApiResponseG[R any] struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	Result      R      `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
}

type TelegramRequest[R, P any] struct {
	method string
	params P
}
type EmptyParams struct{}

func NewRequest[R, P any](method string, params P) TelegramRequest[R, P] {
	return TelegramRequest[R, P]{
		method: method,
		params: params,
	}
}
func (r TelegramRequest[R, P]) Do(bot *Bot) (*R, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r.params)
	if err != nil {
		return nil, err
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

	response := new(ApiResponseG[R])
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
	}
	return &response.Result, nil
}

// request is a low-level call to api.
func (b *Bot) request(methodName string, params any) (map[string]any, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(params)
	if err != nil {
		return nil, err
	}

	if b.debug && b.requestLogger != nil {
		b.requestLogger.Debugln(strings.ReplaceAll(fmt.Sprintf(
			"POST https://api.telegram.org/bot%s/%s %s",
			"<TOKEN>",
			methodName,
			buf.String(),
		), "\n", ""))
	}
	r, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.token, methodName), "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	b.requestLogger.Debugln(fmt.Sprintf("RES %s %s", methodName, string(data)))

	response := new(ApiResponse)
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
	}
	if res, ok := response.Result.(bool); ok {
		return map[string]any{
			"data": res,
		}, nil
	} else if res, ok := response.Result.([]any); ok {
		return map[string]any{
			"data": res,
		}, nil
	} else if res, ok := response.Result.(map[string]any); ok {
		return res, nil
	}
	return map[string]any{}, errors.New("can't parse response")
}
