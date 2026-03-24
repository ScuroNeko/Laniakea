package laniakea

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/slog"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestAutoGenerateCommandsChecksLimitBeforeDelete(t *testing.T) {
	var calls atomic.Int64
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":true}`)),
			}, nil
		}),
	}
	api := tgapi.NewAPI(
		tgapi.NewAPIOpts("token").
			SetAPIUrl("https://example.test").
			SetHTTPClient(client),
	)
	defer func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	plugin := NewPlugin[NoDB]("overflow")
	exec := func(ctx *MsgContext, db *NoDB) {}
	for i := 0; i < 101; i++ {
		plugin.AddCommand(NewCommand(exec, "cmd"+strconv.Itoa(i)))
	}

	bot := &Bot[NoDB]{
		api:     api,
		logger:  slog.CreateLogger(),
		plugins: []Plugin[NoDB]{*plugin},
	}

	err := bot.AutoGenerateCommands()
	if !errors.Is(err, ErrTooManyCommands) {
		t.Fatalf("expected ErrTooManyCommands, got %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("expected no HTTP calls before limit validation, got %d", calls.Load())
	}
}
