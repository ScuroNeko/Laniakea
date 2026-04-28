package laniakea

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

type pollingRoundTripFunc func(*http.Request) (*http.Response, error)

func (f pollingRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type pollingRetryObserver struct {
	recordingObserver
	cancel context.CancelFunc
}

func (o *pollingRetryObserver) OnPollingRetry(ctx context.Context, ev PollingRetryEvent) {
	o.recordingObserver.OnPollingRetry(ctx, ev)
	if o.cancel != nil {
		o.cancel()
	}
}

type testObserver struct{}

func (testObserver) OnReceiveUpdate(context.Context, UpdateReceivedEvent)  {}
func (testObserver) OnHandledUpdate(context.Context, UpdateHandledEvent)   {}
func (testObserver) OnHandlerStarted(context.Context, HandlerStartedEvent) {}
func (testObserver) OnHandlerFinished(context.Context, HandlerFinishedEvent) {
}
func (testObserver) OnSceneTransition(context.Context, SceneTransitionEvent) {}
func (testObserver) OnPolicyChecked(context.Context, PolicyCheckedEvent)     {}
func (testObserver) OnRunnerFinished(context.Context, RunnerFinishedEvent)   {}
func (testObserver) OnPollingRetry(context.Context, PollingRetryEvent)       {}
func (testObserver) OnError(context.Context, ErrorEvent)                     {}

func TestGetUpdateTypesReturnsCopy(t *testing.T) {
	bot := &Bot[NoData]{updateTypes: []tgapi.UpdateType{tgapi.UpdateTypeMessage}}

	got := bot.GetUpdateTypes()
	got[0] = tgapi.UpdateTypeCallbackQuery

	if want := []tgapi.UpdateType{tgapi.UpdateTypeMessage}; !reflect.DeepEqual(bot.updateTypes, want) {
		t.Fatalf("GetUpdateTypes exposed internal slice: got %v want %v", bot.updateTypes, want)
	}
}

func TestAddPluginsSnapshotsConfiguration(t *testing.T) {
	bot := &Bot[NoData]{logger: sneklog.NewLogger()}
	plugin := NewPlugin[NoData]("demo")

	cmd := plugin.NewCommand(func(ctx *MsgContext, db NoData) error { return nil }, "start")
	plugin.AddMiddleware(NewMiddleware("base", func(ctx *MsgContext, db NoData) bool { return true }))

	bot.AddPlugins(plugin)

	cmd.SetDescription("mutated after registration")
	plugin.NewCommand(func(ctx *MsgContext, db NoData) error { return nil }, "late")
	plugin.AddMiddleware(NewMiddleware("late", func(ctx *MsgContext, db NoData) bool { return true }))

	registered := bot.plugins[0]
	if _, exists := registered.commands["late"]; exists {
		t.Fatal("late command leaked into registered plugin snapshot")
	}
	if registered.commands["start"].description != "" {
		t.Fatalf("registered command description unexpectedly mutated: %q", registered.commands["start"].description)
	}
	if len(registered.middlewares) != 1 {
		t.Fatalf("registered middlewares unexpectedly mutated: got %d want 1", len(registered.middlewares))
	}
}

func TestBotPayloadTypeConfiguration(t *testing.T) {
	bot := &Bot[NoData]{payloadType: BotPayloadBase64}

	if got := bot.GetPayloadType(); got != BotPayloadBase64 {
		t.Fatalf("unexpected initial payload type: %q", got)
	}
	bot.SetPayloadType(BotPayloadJSON)
	if got := bot.GetPayloadType(); got != BotPayloadJSON {
		t.Fatalf("unexpected updated payload type: %q", got)
	}
	bot.SetStrictPayloadType(true)
	if !bot.strictPayloadType {
		t.Fatal("expected strict payload type to be enabled")
	}
}

func TestAddPluginsSkipsNilPlugin(t *testing.T) {
	bot := &Bot[NoData]{logger: sneklog.NewLogger()}
	plugin := NewPlugin[NoData]("demo")

	bot.AddPlugins(nil, plugin)

	if len(bot.plugins) != 1 {
		t.Fatalf("expected exactly one registered plugin, got %d", len(bot.plugins))
	}
	if bot.plugins[0].name != "demo" {
		t.Fatalf("unexpected plugin name: %q", bot.plugins[0].name)
	}
}

func TestInitLoggersFallsBackToStdoutLoggerOnFileError(t *testing.T) {
	bot := &Bot[NoData]{}

	bot.initLoggers(&BotOpts{
		Debug:            true,
		WriteToFile:      true,
		UseRequestLogger: true,
		LoggerBasePath:   filepath.Join(t.TempDir(), "missing", "nested"),
	})

	if bot.logger == nil {
		t.Fatal("expected main logger fallback")
	}
	if bot.requestLogger == nil {
		t.Fatal("expected request logger fallback")
	}
	if err := bot.requestLogger.Close(); err != nil {
		t.Fatalf("failed to close request logger: %v", err)
	}
	if err := bot.logger.Close(); err != nil {
		t.Fatalf("failed to close main logger: %v", err)
	}
}

func TestInitLoggersAppliesTokenReplacerToFileLoggers(t *testing.T) {
	tempDir := t.TempDir()
	api := tgapi.NewAPI(tgapi.NewAPIOpts("secret-token"))
	uploader := tgapi.NewUploader(api)
	bot := &Bot[NoData]{token: "secret-token", api: api, uploader: uploader}
	t.Cleanup(func() {
		if err := uploader.Close(); err != nil {
			t.Fatalf("failed to close uploader: %v", err)
		}
	})
	t.Cleanup(func() {
		if err := api.Close(); err != nil {
			t.Fatalf("failed to close api: %v", err)
		}
	})

	bot.initLoggers(&BotOpts{
		Debug:            true,
		WriteToFile:      true,
		UseRequestLogger: true,
		LoggerBasePath:   tempDir,
	})

	apiPath := filepath.Join(tempDir, "api.log")
	apiFile, err := os.OpenFile(apiPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open api log: %v", err)
	}
	defer func() { _ = apiFile.Close() }()
	bot.api.GetLogger().AddWriters(bot.api.GetLogger().CreateTextWriter(apiFile))

	uploaderPath := filepath.Join(tempDir, "uploader.log")
	uploaderFile, err := os.OpenFile(uploaderPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open uploader log: %v", err)
	}
	defer func() { _ = uploaderFile.Close() }()
	bot.uploader.GetLogger().AddWriters(bot.uploader.GetLogger().CreateTextWriter(uploaderFile))

	bot.logger.Infoln("main secret-token")
	bot.requestLogger.Infoln("request secret-token")
	bot.api.GetLogger().Infoln("api secret-token")
	bot.uploader.GetLogger().Infoln("uploader secret-token")

	if err := bot.requestLogger.Close(); err != nil {
		t.Fatalf("failed to close request logger: %v", err)
	}
	if err := bot.logger.Close(); err != nil {
		t.Fatalf("failed to close main logger: %v", err)
	}

	mainLog, err := os.ReadFile(filepath.Join(tempDir, "main.log"))
	if err != nil {
		t.Fatalf("failed to read main log: %v", err)
	}
	requestLog, err := os.ReadFile(filepath.Join(tempDir, "requests.log"))
	if err != nil {
		t.Fatalf("failed to read request log: %v", err)
	}
	apiLog, err := os.ReadFile(apiPath)
	if err != nil {
		t.Fatalf("failed to read api log: %v", err)
	}
	uploaderLog, err := os.ReadFile(uploaderPath)
	if err != nil {
		t.Fatalf("failed to read uploader log: %v", err)
	}

	for _, tt := range []struct {
		name string
		data string
	}{
		{name: "main", data: string(mainLog)},
		{name: "request", data: string(requestLog)},
		{name: "api", data: string(apiLog)},
		{name: "uploader", data: string(uploaderLog)},
	} {
		if strings.Contains(tt.data, "secret-token") {
			t.Fatalf("%s log leaked raw token: %q", tt.name, tt.data)
		}
		if !strings.Contains(tt.data, "<TOKEN>") {
			t.Fatalf("%s log did not contain masked token: %q", tt.name, tt.data)
		}
	}
}

func TestAddPluginsAppliesTokenReplacerToPluginLogger(t *testing.T) {
	bot := &Bot[NoData]{
		token:  "secret-token",
		logger: sneklog.NewLogger(),
	}
	defer func() { _ = bot.logger.Close() }()

	plugin := NewPlugin[NoData]("demo")
	bot.AddPlugins(plugin)

	logPath := filepath.Join(t.TempDir(), "plugin.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open plugin log: %v", err)
	}
	defer func() { _ = file.Close() }()

	bot.plugins[0].logger.AddWriters(bot.plugins[0].logger.CreateTextWriter(file))
	bot.plugins[0].logger.Infoln("plugin secret-token")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read plugin log: %v", err)
	}
	if strings.Contains(string(data), "secret-token") {
		t.Fatalf("plugin log leaked raw token: %q", string(data))
	}
	if !strings.Contains(string(data), "<TOKEN>") {
		t.Fatalf("plugin log did not contain masked token: %q", string(data))
	}
}

func TestNextPollRetryDelay(t *testing.T) {
	tests := []struct {
		name string
		prev time.Duration
		want time.Duration
	}{
		{name: "initial", prev: 0, want: time.Second},
		{name: "double", prev: 2 * time.Second, want: 4 * time.Second},
		{name: "cap", prev: 20 * time.Second, want: 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextPollRetryDelay(tt.prev); got != tt.want {
				t.Fatalf("nextPollRetryDelay(%s) = %s, want %s", tt.prev, got, tt.want)
			}
		})
	}
}

func TestAddDatabaseLoggerWriterSkipsWhenAppDataIsUnset(t *testing.T) {
	bot := &Bot[NoData]{logger: sneklog.NewLogger()}
	called := false

	bot.AddAppDataLoggerWriter(func(db NoData) sneklog.LoggerWriter {
		called = true
		return nil
	})

	if called {
		t.Fatal("expected app-data logger writer to be skipped when app data is unset")
	}
}

func TestAddDatabaseLoggerWriterSkipsWhenAppDataIsNil(t *testing.T) {
	type testDB struct{}

	bot := &Bot[*testDB]{logger: sneklog.NewLogger()}
	var db *testDB
	bot.SetAppData(db)

	called := false
	bot.AddAppDataLoggerWriter(func(db *testDB) sneklog.LoggerWriter {
		called = true
		return nil
	})

	if called {
		t.Fatal("expected app-data logger writer to be skipped when app data is nil")
	}
}

func TestShouldWarnOnValueAppData(t *testing.T) {
	type testDB struct{}
	type dbIface interface{ Ping() error }

	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{name: "NoData", got: shouldWarnOnValueAppData[NoData](), want: false},
		{name: "pointer", got: shouldWarnOnValueAppData[*testDB](), want: false},
		{name: "interface", got: shouldWarnOnValueAppData[dbIface](), want: false},
		{name: "map", got: shouldWarnOnValueAppData[map[string]int](), want: false},
		{name: "struct", got: shouldWarnOnValueAppData[testDB](), want: true},
		{name: "int", got: shouldWarnOnValueAppData[int](), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("shouldWarnOnValueAppData = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestSetAppDataMarksValueWarningOnce(t *testing.T) {
	type testDB struct{}

	bot := &Bot[testDB]{logger: sneklog.NewLogger()}
	bot.SetAppData(testDB{})
	if !bot.warnedValueData {
		t.Fatal("expected value-typed app data to mark warning state")
	}

	ptrBot := &Bot[*testDB]{logger: sneklog.NewLogger()}
	ptrBot.SetAppData(&testDB{})
	if ptrBot.warnedValueData {
		t.Fatal("did not expect pointer-typed app data to mark warning state")
	}
}

func TestSetObserverAndGetObserver(t *testing.T) {
	bot := &Bot[NoData]{logger: sneklog.NewLogger()}
	observer := testObserver{}

	if got := bot.GetObserver(); got != nil {
		t.Fatalf("expected nil observer by default, got %#v", got)
	}

	bot.SetObserver(observer)
	if got := bot.GetObserver(); got == nil {
		t.Fatal("expected observer to be stored")
	}
}

func TestSetObserverNilClearsObserver(t *testing.T) {
	bot := &Bot[NoData]{logger: sneklog.NewLogger()}
	bot.SetObserver(testObserver{})

	if bot.GetObserver() == nil {
		t.Fatal("expected observer to be set")
	}

	bot.SetObserver(nil)
	if got := bot.GetObserver(); got != nil {
		t.Fatalf("expected nil observer after clearing, got %#v", got)
	}
}

func TestRunWithContextRejectsSecondRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		prefixes:    []string{"/"},
		plugins:     []Plugin[NoData]{{name: "demo"}},
		updateQueue: make(chan *tgapi.Update, 1),
		maxWorkers:  1,
	}

	if err := bot.RunWithContext(ctx); err != nil {
		t.Fatalf("first RunWithContext returned error: %v", err)
	}
	if err := bot.RunWithContext(ctx); !errors.Is(err, ErrBotAlreadyRun) {
		t.Fatalf("expected ErrBotAlreadyRun on second run, got %v", err)
	}
}

func TestRunWithContextKeepsEnabledRequestLogger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	requestLogger := sneklog.NewLogger()
	bot := &Bot[NoData]{
		logger:        sneklog.NewLogger(),
		requestLogger: requestLogger,
		useReqLogger:  true,
		prefixes:      []string{"/"},
		plugins:       []Plugin[NoData]{{name: "demo"}},
		updateQueue:   make(chan *tgapi.Update, 1),
		maxWorkers:    1,
	}
	t.Cleanup(func() {
		_ = bot.Close()
	})

	if err := bot.RunWithContext(ctx); err != nil {
		t.Fatalf("RunWithContext returned error: %v", err)
	}
	if got := bot.GetRequestLogger(); got != requestLogger {
		t.Fatalf("expected enabled request logger to be preserved, got %#v", got)
	}
}

func TestCloseDoesNotDeleteWebhook(t *testing.T) {
	requests := 0
	client := &http.Client{
		Transport: pollingRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":true}`)),
			}, nil
		}),
	}

	api := tgapi.NewAPI(
		tgapi.NewAPIOpts("token").
			SetAPIURL("http://example.invalid").
			SetHTTPClient(client),
	)
	uploader := tgapi.NewUploader(api)

	bot := &Bot[NoData]{
		logger:        sneklog.NewLogger(),
		webHookLogger: sneklog.NewLogger(),
		api:           api,
		uploader:      uploader,
	}

	if err := bot.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if requests != 0 {
		t.Fatalf("Close performed unexpected remote requests: got %d want 0", requests)
	}
}

func TestRunWithContextEmitsPollingRetryAndErrorEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observer := &pollingRetryObserver{cancel: cancel}

	client := &http.Client{
		Transport: pollingRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":false,"error_code":500,"description":"boom"}`)),
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
		logger:      sneklog.NewLogger(),
		api:         api,
		prefixes:    []string{"/"},
		plugins:     []Plugin[NoData]{{name: "demo"}},
		updateQueue: make(chan *tgapi.Update, 1),
		maxWorkers:  1,
		observer:    observer,
	}

	if err := bot.RunWithContext(ctx); err != nil {
		t.Fatalf("RunWithContext returned error: %v", err)
	}

	if len(observer.retries) != 1 {
		t.Fatalf("expected one polling retry event, got %d", len(observer.retries))
	}
	if got := observer.retries[0]; got.Attempt != 1 || got.Delay <= 0 || got.Err == nil {
		t.Fatalf("unexpected polling retry event: %#v", got)
	}
	if len(observer.errors) != 1 {
		t.Fatalf("expected one polling error event, got %d", len(observer.errors))
	}
	if got := observer.errors[0]; got.HandlerKind != HandlerPollingKind || got.HandlerName != "getUpdates" || got.Plugin != "bot" || got.Err == nil || got.UserFacing {
		t.Fatalf("unexpected polling error event: %#v", got)
	}
}

func TestBotConfigurationFreezesAfterRunStarts(t *testing.T) {
	type testDB struct{ Name string }

	makeBot := func() *Bot[*testDB] {
		return &Bot[*testDB]{
			logger:             sneklog.NewLogger(),
			prefixes:           []string{"/"},
			updateTypes:        []tgapi.UpdateType{tgapi.UpdateTypeMessage},
			payloadType:        BotPayloadBase64,
			strictPayloadType:  false,
			errorTemplate:      "%s",
			l10n:               &L10n{},
			draftProvider:      &DraftProvider{},
			sessionStore:       NewMemorySessionStore(),
			sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		}
	}

	tests := []struct {
		name  string
		check func(t *testing.T, bot *Bot[*testDB])
	}{
		{
			name: "SetAppData",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := &testDB{Name: "before"}
				bot.SetAppData(original)
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				later := &testDB{Name: "after"}
				bot.SetAppData(later)
				if bot.appData != original {
					t.Fatal("SetAppData mutated after configuration freeze")
				}
			},
		},
		{
			name: "UpdateTypes",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := append([]tgapi.UpdateType(nil), bot.updateTypes...)
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetUpdateTypes(tgapi.UpdateTypePoll)
				if !reflect.DeepEqual(bot.updateTypes, original) {
					t.Fatalf("UpdateTypes mutated after configuration freeze: got %v want %v", bot.updateTypes, original)
				}
			},
		},
		{
			name: "AddUpdateType",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := append([]tgapi.UpdateType(nil), bot.updateTypes...)
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.AddUpdateType(tgapi.UpdateTypePoll)
				if !reflect.DeepEqual(bot.updateTypes, original) {
					t.Fatalf("AddUpdateType mutated after configuration freeze: got %v want %v", bot.updateTypes, original)
				}
			},
		},
		{
			name: "SetPayloadType",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetPayloadType(BotPayloadJSON)
				if bot.payloadType != BotPayloadBase64 {
					t.Fatalf("payloadType mutated after configuration freeze: got %q want %q", bot.payloadType, BotPayloadBase64)
				}
			},
		},
		{
			name: "SetStrictPayloadType",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetStrictPayloadType(true)
				if bot.strictPayloadType {
					t.Fatal("strictPayloadType mutated after configuration freeze")
				}
			},
		},
		{
			name: "AddPrefixes",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := append([]string(nil), bot.prefixes...)
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.AddPrefixes("!")
				if !reflect.DeepEqual(bot.prefixes, original) {
					t.Fatalf("prefixes mutated after configuration freeze: got %v want %v", bot.prefixes, original)
				}
			},
		},
		{
			name: "ErrorTemplate",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetErrorTemplate("changed")
				if bot.errorTemplate != "%s" {
					t.Fatalf("errorTemplate mutated after configuration freeze: got %q want %q", bot.errorTemplate, "%s")
				}
			},
		},
		{
			name: "SetDraftProvider",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := bot.draftProvider
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetDraftProvider(&DraftProvider{})
				if bot.draftProvider != original {
					t.Fatal("draftProvider mutated after configuration freeze")
				}
			},
		},
		{
			name: "SetSessionStore",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := bot.sessionStore
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetSessionStore(NewMemorySessionStore())
				if bot.sessionStore != original {
					t.Fatal("sessionStore mutated after configuration freeze")
				}
			},
		},
		{
			name: "SetSceneScopePriority",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := append([]SceneScope(nil), bot.sceneScopePriority...)
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetSceneScopePriority([]SceneScope{SceneScopeUser})
				if !reflect.DeepEqual(bot.sceneScopePriority, original) {
					t.Fatalf("sceneScopePriority mutated after configuration freeze: got %v want %v", bot.sceneScopePriority, original)
				}
			},
		},
		{
			name: "AddL10n",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := bot.l10n
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				bot.SetL10n(&L10n{})
				if bot.l10n != original {
					t.Fatal("l10n mutated after configuration freeze")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, makeBot())
		})
	}
}

func TestAddPluginsAndRuntimeRegistrationsNoOpAfterRunStarts(t *testing.T) {
	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		prefixes:    []string{"/"},
		middlewares: []Middleware[NoData]{NewMiddleware("base", func(ctx *MsgContext, db NoData) bool { return true })},
		runners:     []Runner[NoData]{NewRunner("base", func(bot *Bot[NoData]) error { return nil })},
	}
	plugin := NewPlugin[NoData]("late")

	if err := bot.beginRun(); err != nil {
		t.Fatalf("beginRun returned error: %v", err)
	}
	defer bot.finishRun()

	bot.AddPlugins(plugin)
	bot.AddMiddleware(NewMiddleware("late", func(ctx *MsgContext, db NoData) bool { return true }))
	bot.AddRunner(NewRunner("late", func(bot *Bot[NoData]) error { return nil }))

	if len(bot.plugins) != 0 {
		t.Fatalf("expected AddPlugins to be ignored after configuration freeze, got %d plugins", len(bot.plugins))
	}
	if len(bot.middlewares) != 1 {
		t.Fatalf("expected AddMiddleware to be ignored after configuration freeze, got %d middlewares", len(bot.middlewares))
	}
	if len(bot.runners) != 1 {
		t.Fatalf("expected AddRunner to be ignored after configuration freeze, got %d runners", len(bot.runners))
	}
}
