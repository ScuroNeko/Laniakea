package laniakea

import (
	"context"
	"errors"
	"testing"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/slog"
)

func TestPluginAddSceneRegistersScene(t *testing.T) {
	plugin := NewPlugin[NoDB]("wizard")
	scene := NewScene[NoDB]("signup")

	plugin.AddScene(scene)

	if got, ok := plugin.scenes["signup"]; !ok || got != scene {
		t.Fatalf("scene was not registered in plugin: ok=%v got=%p want=%p", ok, got, scene)
	}
	if scene.PluginName != "wizard" {
		t.Fatalf("unexpected plugin name on scene: got %q want %q", scene.PluginName, "wizard")
	}
}

func TestBotAddPluginsPreservesScenesAndHandlesThem(t *testing.T) {
	called := false

	plugin := NewPlugin[NoDB]("wizard")
	plugin.NewScene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoDB) (SceneResult, error) {
			called = true
			if ctx.Text != "hello there" {
				t.Fatalf("unexpected scene text: got %q want %q", ctx.Text, "hello there")
			}
			return ctx.Exit(), nil
		})

	bot := &Bot[NoDB]{
		logger:             slog.CreateLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
	}
	bot.AddPlugins(plugin)

	sceneMeta, ok := bot.findScene("signup")
	if !ok {
		t.Fatal("expected scene metadata to be available after plugin registration")
	}
	if sceneMeta.Entry != "start" {
		t.Fatalf("unexpected scene entry: got %q want %q", sceneMeta.Entry, "start")
	}

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: string(tgapi.ChatTypePrivate)}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 1,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 7,
			Text:      "hello there",
			Chat:      &tgapi.Chat{ID: 100, Type: string(tgapi.ChatTypePrivate)},
			From:      &tgapi.User{ID: 42},
		},
	})

	if !called {
		t.Fatal("expected scene step handler to be called")
	}

	lookupCtx := &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: string(tgapi.ChatTypePrivate)}},
		FromID: 42,
	}
	if _, session, err := bot.findSceneSession(lookupCtx); err == nil && session.Scene != "" {
		t.Fatalf("expected scene session to be removed after exit, got %#v", session)
	}
}

func TestBuildSceneKeyRejectsMissingContextFields(t *testing.T) {
	tests := []struct {
		name  string
		scope SceneScope
		ctx   *MsgContext
	}{
		{
			name:  "nil context",
			scope: SceneScopeUserChat,
			ctx:   nil,
		},
		{
			name:  "missing message for chat scope",
			scope: SceneScopeChat,
			ctx:   &MsgContext{},
		},
		{
			name:  "missing from id for user scope",
			scope: SceneScopeUser,
			ctx:   &MsgContext{},
		},
		{
			name:  "missing from id for user chat scope",
			scope: SceneScopeUserChat,
			ctx: &MsgContext{
				Msg: &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: string(tgapi.ChatTypePrivate)}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if key, ok := buildSceneKey(tt.scope, tt.ctx); ok || key != "" {
				t.Fatalf("expected invalid scene key, got key=%q ok=%v", key, ok)
			}
		})
	}
}

func TestEnterSceneRejectsMissingEntryConfiguration(t *testing.T) {
	t.Run("empty entry", func(t *testing.T) {
		plugin := NewPlugin[NoDB]("wizard")
		plugin.NewScene("signup")

		bot := &Bot[NoDB]{
			logger:             slog.CreateLogger(),
			sessionStore:       NewMemorySessionStore(),
			sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		}
		bot.AddPlugins(plugin)

		ctx := &MsgContext{
			Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: string(tgapi.ChatTypePrivate)}},
			FromID:       42,
			sceneRuntime: bot,
		}

		err := ctx.EnterScene("signup")
		if !errors.Is(err, ErrSceneEntryNotSet) {
			t.Fatalf("expected ErrSceneEntryNotSet, got %v", err)
		}
	})

	t.Run("missing entry step", func(t *testing.T) {
		plugin := NewPlugin[NoDB]("wizard")
		plugin.NewScene("signup").SetEntry("start")

		bot := &Bot[NoDB]{
			logger:             slog.CreateLogger(),
			sessionStore:       NewMemorySessionStore(),
			sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		}
		bot.AddPlugins(plugin)

		ctx := &MsgContext{
			Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: string(tgapi.ChatTypePrivate)}},
			FromID:       42,
			sceneRuntime: bot,
		}

		err := ctx.EnterScene("signup")
		if !errors.Is(err, ErrSceneStepNotFound) {
			t.Fatalf("expected ErrSceneStepNotFound, got %v", err)
		}
	})
}
