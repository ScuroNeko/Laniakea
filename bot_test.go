package laniakea

import (
	"path/filepath"
	"reflect"
	"testing"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/slog"
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

	cmd := plugin.NewCommand(func(ctx *MsgContext, db *NoDB) {}, "start")
	plugin.AddMiddleware(*NewMiddleware("base", func(ctx *MsgContext, db *NoDB) bool { return true }))

	bot.AddPlugins(plugin)

	cmd.SetDescription("mutated after registration")
	plugin.NewCommand(func(ctx *MsgContext, db *NoDB) {}, "late")
	plugin.AddMiddleware(*NewMiddleware("late", func(ctx *MsgContext, db *NoDB) bool { return true }))

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
