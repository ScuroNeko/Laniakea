package laniakea

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
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

	plugin := NewPlugin[NoData]("overflow")
	exec := func(ctx *MsgContext, db NoData) error { return nil }
	for i := 0; i < 101; i++ {
		plugin.AddCommand(NewCommand(exec, "cmd"+strconv.Itoa(i)))
	}

	bot := &Bot[NoData]{
		api:     api,
		logger:  slog.CreateLogger(),
		plugins: []Plugin[NoData]{*plugin},
	}

	err := bot.AutoGenerateCommands()
	if !errors.Is(err, ErrTooManyCommands) {
		t.Fatalf("expected ErrTooManyCommands, got %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("expected no HTTP calls before limit validation, got %d", calls.Load())
	}
}

func TestGatherCommandsForPluginReturnsSortedCommands(t *testing.T) {
	plugin := NewPlugin[NoData]("sorted")
	exec := func(ctx *MsgContext, db NoData) error { return nil }

	plugin.AddCommand(NewCommand(exec, "zeta"))
	plugin.AddCommand(NewCommand(exec, "alpha"))
	plugin.AddCommand(NewCommand(exec, "mid"))

	commands := gatherCommandsForPlugin(*plugin)
	got := make([]string, 0, len(commands))
	for _, cmd := range commands {
		got = append(got, cmd.Command)
	}

	want := []string{"alpha", "mid", "zeta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected command order: got %v want %v", got, want)
	}
}
