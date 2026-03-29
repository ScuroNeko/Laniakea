package laniakea

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
)

func TestGetUpdateTypesReturnsCopy(t *testing.T) {
	bot := &Bot[NoDB]{updateTypes: []tgapi.UpdateType{tgapi.UpdateTypeMessage}}

	got := bot.GetUpdateTypes()
	got[0] = tgapi.UpdateTypeCallbackQuery

	if want := []tgapi.UpdateType{tgapi.UpdateTypeMessage}; !reflect.DeepEqual(bot.updateTypes, want) {
		t.Fatalf("GetUpdateTypes exposed internal slice: got %v want %v", bot.updateTypes, want)
	}
}

func TestAddPluginsSnapshotsConfiguration(t *testing.T) {
	bot := &Bot[NoDB]{logger: slog.CreateLogger()}
	plugin := NewPlugin[NoDB]("demo")

	cmd := plugin.NewCommand(func(ctx *MsgContext, db NoDB) error { return nil }, "start")
	plugin.AddMiddleware(NewMiddleware("base", func(ctx *MsgContext, db NoDB) bool { return true }))

	bot.AddPlugins(plugin)

	cmd.SetDescription("mutated after registration")
	plugin.NewCommand(func(ctx *MsgContext, db NoDB) error { return nil }, "late")
	plugin.AddMiddleware(NewMiddleware("late", func(ctx *MsgContext, db NoDB) bool { return true }))

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
	bot := &Bot[NoDB]{payloadType: BotPayloadBase64}

	if got := bot.GetPayloadType(); got != BotPayloadBase64 {
		t.Fatalf("unexpected initial payload type: %q", got)
	}
	bot.SetPayloadType(BotPayloadJson)
	if got := bot.GetPayloadType(); got != BotPayloadJson {
		t.Fatalf("unexpected updated payload type: %q", got)
	}
	bot.SetStrictPayloadType(true)
	if !bot.strictPayloadType {
		t.Fatal("expected strict payload type to be enabled")
	}
}

func TestAddPluginsSkipsNilPlugin(t *testing.T) {
	bot := &Bot[NoDB]{logger: slog.CreateLogger()}
	plugin := NewPlugin[NoDB]("demo")

	bot.AddPlugins(nil, plugin)

	if len(bot.plugins) != 1 {
		t.Fatalf("expected exactly one registered plugin, got %d", len(bot.plugins))
	}
	if bot.plugins[0].name != "demo" {
		t.Fatalf("unexpected plugin name: %q", bot.plugins[0].name)
	}
}

func TestInitLoggersFallsBackToStdoutLoggerOnFileError(t *testing.T) {
	bot := &Bot[NoDB]{}

	bot.initLoggers(&BotOpts{
		Debug:            true,
		WriteToFile:      true,
		UseRequestLogger: true,
		LoggerBasePath:   filepath.Join(t.TempDir(), "missing", "nested"),
	})

	if bot.logger == nil {
		t.Fatal("expected main logger fallback")
	}
	if bot.RequestLogger == nil {
		t.Fatal("expected request logger fallback")
	}
	if err := bot.RequestLogger.Close(); err != nil {
		t.Fatalf("failed to close request logger: %v", err)
	}
	if err := bot.logger.Close(); err != nil {
		t.Fatalf("failed to close main logger: %v", err)
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

func TestAddDatabaseLoggerWriterSkipsWhenDBContextIsUnset(t *testing.T) {
	bot := &Bot[NoDB]{logger: slog.CreateLogger()}
	called := false

	bot.AddDatabaseLoggerWriter(func(db NoDB) slog.LoggerWriter {
		called = true
		return nil
	})

	if called {
		t.Fatal("expected database logger writer to be skipped when db context is unset")
	}
}

func TestAddDatabaseLoggerWriterSkipsWhenDBContextIsNil(t *testing.T) {
	type testDB struct{}

	bot := &Bot[*testDB]{logger: slog.CreateLogger()}
	var db *testDB
	bot.DatabaseContext(db)

	called := false
	bot.AddDatabaseLoggerWriter(func(db *testDB) slog.LoggerWriter {
		called = true
		return nil
	})

	if called {
		t.Fatal("expected database logger writer to be skipped when db context is nil")
	}
}

func TestShouldWarnOnValueDBContext(t *testing.T) {
	type testDB struct{}
	type dbIface interface{ Ping() error }

	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{name: "NoDB", got: shouldWarnOnValueDBContext[NoDB](), want: false},
		{name: "pointer", got: shouldWarnOnValueDBContext[*testDB](), want: false},
		{name: "interface", got: shouldWarnOnValueDBContext[dbIface](), want: false},
		{name: "map", got: shouldWarnOnValueDBContext[map[string]int](), want: false},
		{name: "struct", got: shouldWarnOnValueDBContext[testDB](), want: true},
		{name: "int", got: shouldWarnOnValueDBContext[int](), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("shouldWarnOnValueDBContext = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestDatabaseContextMarksValueWarningOnce(t *testing.T) {
	type testDB struct{}

	bot := &Bot[testDB]{logger: slog.CreateLogger()}
	bot.DatabaseContext(testDB{})
	if !bot.warnedValueDB {
		t.Fatal("expected value-typed database context to mark warning state")
	}

	ptrBot := &Bot[*testDB]{logger: slog.CreateLogger()}
	ptrBot.DatabaseContext(&testDB{})
	if ptrBot.warnedValueDB {
		t.Fatal("did not expect pointer-typed database context to mark warning state")
	}
}

func TestRunWithContextRejectsSecondRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	bot := &Bot[NoDB]{
		logger:      slog.CreateLogger(),
		prefixes:    []string{"/"},
		plugins:     []Plugin[NoDB]{{name: "demo"}},
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

func TestBotConfigurationFreezesAfterRunStarts(t *testing.T) {
	type testDB struct{ Name string }

	makeBot := func() *Bot[*testDB] {
		return &Bot[*testDB]{
			logger:             slog.CreateLogger(),
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
		name string
		check func(t *testing.T, bot *Bot[*testDB])
	}{
		{
			name: "DatabaseContext",
			check: func(t *testing.T, bot *Bot[*testDB]) {
				original := &testDB{Name: "before"}
				bot.DatabaseContext(original)
				if err := bot.beginRun(); err != nil {
					t.Fatalf("beginRun returned error: %v", err)
				}
				t.Cleanup(bot.finishRun)

				later := &testDB{Name: "after"}
				bot.DatabaseContext(later)
				if bot.dbContext != original {
					t.Fatal("DatabaseContext mutated after configuration freeze")
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

				bot.UpdateTypes(tgapi.UpdateTypePoll)
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

				bot.SetPayloadType(BotPayloadJson)
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

				bot.ErrorTemplate("changed")
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

				bot.AddL10n(&L10n{})
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
	bot := &Bot[NoDB]{
		logger:      slog.CreateLogger(),
		prefixes:    []string{"/"},
		middlewares: []Middleware[NoDB]{NewMiddleware("base", func(ctx *MsgContext, db NoDB) bool { return true })},
		runners:     []Runner[NoDB]{NewRunner("base", func(bot *Bot[NoDB]) error { return nil })},
	}
	plugin := NewPlugin[NoDB]("late")

	if err := bot.beginRun(); err != nil {
		t.Fatalf("beginRun returned error: %v", err)
	}
	defer bot.finishRun()

	bot.AddPlugins(plugin)
	bot.AddMiddleware(NewMiddleware("late", func(ctx *MsgContext, db NoDB) bool { return true }))
	bot.AddRunner(NewRunner("late", func(bot *Bot[NoDB]) error { return nil }))

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
