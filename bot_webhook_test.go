package laniakea

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

func TestEnqueueUpdateCopiesValue(t *testing.T) {
	bot := &Bot[NoData]{
		updateQueue: make(chan *tgapi.Update, 1),
	}

	update := tgapi.Update{UpdateID: 42}
	if err := bot.enqueueUpdate(context.Background(), update); err != nil {
		t.Fatalf("enqueueUpdate returned error: %v", err)
	}

	update.UpdateID = 99

	got := <-bot.updateQueue
	if got.UpdateID != 42 {
		t.Fatalf("enqueueUpdate did not copy the update value: got %d want %d", got.UpdateID, 42)
	}
}

func TestUpdateHandlerEnqueuesUpdate(t *testing.T) {
	bot := &Bot[NoData]{
		updateQueue:   make(chan *tgapi.Update, 1),
		webhookLogger: sneklog.NewLogger(),
	}
	t.Cleanup(func() {
		_ = bot.webhookLogger.Close()
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"update_id":7,"message":{"message_id":1,"date":1,"chat":{"id":1,"type":"private"},"text":"/start"}}`))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "secret")
	rec := httptest.NewRecorder()

	updateHandler(context.Background(), bot, []byte("secret")).ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Result().StatusCode, http.StatusOK)
	}

	select {
	case got := <-bot.updateQueue:
		if got.UpdateID != 7 {
			t.Fatalf("unexpected update id in queue: got %d want %d", got.UpdateID, 7)
		}
		if got.Type != tgapi.UpdateTypeMessage {
			t.Fatalf("unexpected update type in queue: got %q want %q", got.Type, tgapi.UpdateTypeMessage)
		}
	default:
		t.Fatal("expected webhook handler to enqueue an update")
	}
}

func TestRunWebhookRuntimeRejectsSecondRun(t *testing.T) {
	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		updateQueue: make(chan *tgapi.Update, 1),
		maxWorkers:  1,
	}
	t.Cleanup(func() {
		_ = bot.logger.Close()
	})

	if err := bot.runWebhookRuntime(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatalf("first runWebhookRuntime returned error: %v", err)
	}
	if err := bot.runWebhookRuntime(context.Background(), func(context.Context) error { return nil }); !errors.Is(err, ErrBotAlreadyRun) {
		t.Fatalf("expected ErrBotAlreadyRun on second run, got %v", err)
	}
}

func TestRunWebhookRuntimeExecutesRunners(t *testing.T) {
	var calls atomic.Int32

	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		updateQueue: make(chan *tgapi.Update, 1),
		maxWorkers:  1,
		runners: []Runner[NoData]{
			NewRunner("runner", func(bot *Bot[NoData]) error {
				calls.Add(1)
				return nil
			}).Async(false),
		},
	}
	t.Cleanup(func() {
		_ = bot.logger.Close()
	})

	if err := bot.runWebhookRuntime(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatalf("runWebhookRuntime returned error: %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected runner to execute once, got %d", got)
	}
}

func TestRunWebhookRuntimePreservesConfiguredWebhookLogger(t *testing.T) {
	webhookLogger := sneklog.NewLogger()
	bot := &Bot[NoData]{
		logger:        sneklog.NewLogger(),
		webhookLogger: webhookLogger,
		updateQueue:   make(chan *tgapi.Update, 1),
		maxWorkers:    1,
	}
	t.Cleanup(func() {
		_ = bot.logger.Close()
		if bot.webhookLogger != nil {
			_ = bot.webhookLogger.Close()
		}
	})

	if err := bot.runWebhookRuntime(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatalf("runWebhookRuntime returned error: %v", err)
	}
	if bot.webhookLogger != webhookLogger {
		t.Fatal("expected runWebhookRuntime to preserve configured webhook logger")
	}
}

func TestRunWebhookRuntimeProcessesEnqueuedUpdate(t *testing.T) {
	var calls atomic.Int32
	plugin := NewPlugin[NoData]("demo")
	plugin.Command("start", func(ctx *MessageContext, db NoData) error {
		calls.Add(1)
		return nil
	})

	bot := &Bot[NoData]{
		logger:        sneklog.NewLogger(),
		webhookLogger: sneklog.NewLogger(),
		prefixes:      []string{"/"},
		plugins:       []Plugin[NoData]{*plugin},
		updateQueue:   make(chan *tgapi.Update, 1),
		maxWorkers:    1,
	}
	t.Cleanup(func() {
		_ = bot.logger.Close()
		_ = bot.webhookLogger.Close()
	})

	err := bot.runWebhookRuntime(context.Background(), func(ctx context.Context) error {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"update_id":9,"message":{"message_id":1,"date":1,"chat":{"id":1,"type":"private"},"from":{"id":2,"is_bot":false,"first_name":"Test"},"text":"/start"}}`))
		rec := httptest.NewRecorder()

		updateHandler(ctx, bot, []byte("")).ServeHTTP(rec, req)
		if rec.Result().StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: got %d want %d", rec.Result().StatusCode, http.StatusOK)
		}

		deadline := time.After(time.Second)
		for calls.Load() == 0 {
			select {
			case <-deadline:
				t.Fatal("webhook runtime did not process enqueued update")
			default:
				time.Sleep(time.Millisecond)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("runWebhookRuntime returned error: %v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected command handler to run once, got %d", calls.Load())
	}
}

func TestWebhookAllowedUpdatesUsesBotUpdateTypesByDefault(t *testing.T) {
	bot := &Bot[NoData]{
		updateTypes: []tgapi.UpdateType{tgapi.UpdateTypeMessage, tgapi.UpdateTypeCallbackQuery},
	}
	opts := NewBotWebhookOpts()

	got := bot.webhookAllowedUpdates(opts)
	if len(got) != 2 {
		t.Fatalf("unexpected allowed updates length: got %d want %d", len(got), 2)
	}
	if got[0] != tgapi.UpdateTypeMessage || got[1] != tgapi.UpdateTypeCallbackQuery {
		t.Fatalf("unexpected allowed updates: %v", got)
	}

	got[0] = tgapi.UpdateTypePoll
	if bot.updateTypes[0] != tgapi.UpdateTypeMessage {
		t.Fatalf("webhookAllowedUpdates exposed internal slice: got %v", bot.updateTypes)
	}
}

func TestValidateWebhookPath(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		useStatusPath bool
		wantErr       bool
	}{
		{name: "root", path: "/", wantErr: false},
		{name: "custom path", path: "/telegram", wantErr: false},
		{name: "empty", path: "", wantErr: true},
		{name: "missing slash", path: "telegram", wantErr: true},
		{name: "query", path: "/telegram?x=1", wantErr: true},
		{name: "fragment", path: "/telegram#main", wantErr: true},
		{name: "status collision", path: "/status", useStatusPath: true, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateWebhookPath(tc.path, tc.useStatusPath)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateWebhookTLSFiles(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		wantErr bool
	}{
		{name: "no tls", files: nil, wantErr: false},
		{name: "two files", files: []string{"key.pem", "cert.pem"}, wantErr: false},
		{name: "one file", files: []string{"cert.pem"}, wantErr: true},
		{name: "three files", files: []string{"a", "b", "c"}, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateWebhookTLSFiles(tc.files)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUpdateHandlerRejectsOversizedBody(t *testing.T) {
	bot := &Bot[NoData]{
		updateQueue:   make(chan *tgapi.Update, 1),
		webhookLogger: sneklog.NewLogger(),
	}
	t.Cleanup(func() {
		_ = bot.webhookLogger.Close()
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(strings.Repeat("a", (256<<10)+1)))
	rec := httptest.NewRecorder()

	updateHandler(context.Background(), bot, []byte("")).ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("unexpected status: got %d want %d", rec.Result().StatusCode, http.StatusRequestEntityTooLarge)
	}
}

func TestStatusHandlerRequiresMatchingSecret(t *testing.T) {
	client := &http.Client{
		Transport: pollingRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"url":"https://bot.example.com/telegram"}}`)),
			}, nil
		}),
	}
	api := tgapi.NewAPI(
		tgapi.NewAPIOpts("token").
			SetAPIURL("http://example.invalid").
			SetHTTPClient(client),
	)
	defer func() {
		_ = api.Close()
	}()

	bot := &Bot[NoData]{
		api:           api,
		webhookLogger: sneklog.NewLogger(),
	}
	t.Cleanup(func() {
		_ = bot.webhookLogger.Close()
	})

	handler := statusHandler(bot, []byte("secret"))

	tests := []struct {
		name       string
		headerName string
		headerVal  string
		wantStatus int
	}{
		{name: "missing auth", wantStatus: http.StatusNotFound},
		{name: "wrong auth", headerName: "Authorization", headerVal: "wrong", wantStatus: http.StatusNotFound},
		{name: "matching telegram header", headerName: "X-Telegram-Bot-Api-Secret-Token", headerVal: "secret", wantStatus: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/status", nil)
			if tc.headerName != "" {
				req.Header.Set(tc.headerName, tc.headerVal)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Result().StatusCode != tc.wantStatus {
				t.Fatalf("unexpected status: got %d want %d", rec.Result().StatusCode, tc.wantStatus)
			}
		})
	}
}

func TestRunWebhookWithContextRejectsInvalidTLSFilesBeforeRemoteSetup(t *testing.T) {
	bot := &Bot[NoData]{
		prefixes: []string{"/"},
		plugins:  []Plugin[NoData]{{name: "demo"}},
	}
	opts := NewBotWebhookOpts().SetURL("https://bot.example.com")

	err := bot.RunWebhookWithContext(context.Background(), opts, "cert.pem")
	if err == nil {
		t.Fatal("expected tls validation error, got nil")
	}
	if !strings.Contains(err.Error(), "both private and public keys") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWebhookWithContextRequiresSecretWhenStatusPathEnabled(t *testing.T) {
	bot := &Bot[NoData]{
		prefixes: []string{"/"},
		plugins:  []Plugin[NoData]{{name: "demo"}},
	}
	opts := NewBotWebhookOpts().
		SetURL("https://bot.example.com").
		SetUseStatusPath(true)

	err := bot.RunWebhookWithContext(context.Background(), opts)
	if err == nil {
		t.Fatal("expected status-path secret validation error, got nil")
	}
	if !strings.Contains(err.Error(), "SecretToken required") {
		t.Fatalf("unexpected error: %v", err)
	}
}
