package laniakea

import (
	"context"
	"errors"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

type failingSessionStore struct {
	getErr    error
	setErr    error
	deleteErr error
}

func (s failingSessionStore) Get(string) (SceneSession, error) {
	return SceneSession{}, s.getErr
}

func (s failingSessionStore) Set(string, SceneSession) error {
	return s.setErr
}

func (s failingSessionStore) Delete(string) error {
	return s.deleteErr
}

func TestPluginAddSceneRegistersScene(t *testing.T) {
	plugin := NewPlugin[NoData]("wizard")
	scene := NewScene[NoData]("signup")

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

	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			called = true
			if ctx.Text != "hello there" {
				t.Fatalf("unexpected scene text: got %q want %q", ctx.Text, "hello there")
			}
			return ctx.Exit(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
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
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
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
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if !called {
		t.Fatal("expected scene step handler to be called")
	}

	lookupCtx := &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
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
				Msg: &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
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
		plugin := NewPlugin[NoData]("wizard")
		plugin.Scene("signup")

		bot := &Bot[NoData]{
			logger:             sneklog.NewLogger(),
			sessionStore:       NewMemorySessionStore(),
			sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		}
		bot.AddPlugins(plugin)

		ctx := &MsgContext{
			Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
			FromID:       42,
			sceneRuntime: bot,
		}

		err := ctx.EnterScene("signup")
		if !errors.Is(err, ErrSceneEntryNotSet) {
			t.Fatalf("expected ErrSceneEntryNotSet, got %v", err)
		}
	})

	t.Run("missing entry step", func(t *testing.T) {
		plugin := NewPlugin[NoData]("wizard")
		plugin.Scene("signup").SetEntry("start")

		bot := &Bot[NoData]{
			logger:             sneklog.NewLogger(),
			sessionStore:       NewMemorySessionStore(),
			sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		}
		bot.AddPlugins(plugin)

		ctx := &MsgContext{
			Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
			FromID:       42,
			sceneRuntime: bot,
		}

		err := ctx.EnterScene("signup")
		if !errors.Is(err, ErrSceneStepNotFound) {
			t.Fatalf("expected ErrSceneStepNotFound, got %v", err)
		}
	})
}

func TestSceneContextMethodsRequireRuntime(t *testing.T) {
	ctx := &MsgContext{}

	if err := ctx.EnterScene("signup"); !errors.Is(err, ErrSceneRuntimeNil) {
		t.Fatalf("expected ErrSceneRuntimeNil from EnterScene, got %v", err)
	}
	if err := ctx.EnterSceneStep("signup", "start"); !errors.Is(err, ErrSceneRuntimeNil) {
		t.Fatalf("expected ErrSceneRuntimeNil from EnterSceneStep, got %v", err)
	}
	if err := ctx.ExitScene(); !errors.Is(err, ErrSceneRuntimeNil) {
		t.Fatalf("expected ErrSceneRuntimeNil from ExitScene, got %v", err)
	}
}

func TestSceneCommandHandlerRunsBeforeStep(t *testing.T) {
	sceneCommandCalled := false
	stepCalled := false

	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			stepCalled = true
			return ctx.Stay(), nil
		}).
		OnCommand("cancel", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			sceneCommandCalled = true
			if ctx.Prefix != "/" {
				t.Fatalf("unexpected prefix: got %q want /", ctx.Prefix)
			}
			if ctx.Text != "right now" {
				t.Fatalf("unexpected scene command text: got %q want %q", ctx.Text, "right now")
			}
			if len(ctx.Args) != 2 || ctx.Args[0] != "right" || ctx.Args[1] != "now" {
				t.Fatalf("unexpected scene command args: %#v", ctx.Args)
			}
			return ctx.Exit(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 2,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 8,
			Text:      "/cancel right now",
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if !sceneCommandCalled {
		t.Fatal("expected scene command handler to be called")
	}
	if stepCalled {
		t.Fatal("expected scene command to short-circuit the scene step")
	}
}

func TestSceneCommandObserverEmitsLifecycleEvents(t *testing.T) {
	observer := &recordingObserver{}
	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Stay(), nil
		}).
		OnCommand("cancel", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Exit(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		observer:           observer,
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 22,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 9,
			Text:      "/cancel",
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if len(observer.started) != 1 {
		t.Fatalf("expected one scene started event, got %d", len(observer.started))
	}
	if got := observer.started[0]; got.HandlerKind != HandlerSceneCommandKind || got.HandlerName != "cancel" || got.Plugin != "wizard" {
		t.Fatalf("unexpected scene started event: %#v", got)
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one scene finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; got.HandlerKind != HandlerSceneCommandKind || got.HandlerName != "cancel" || got.Plugin != "wizard" || got.Err != nil {
		t.Fatalf("unexpected scene finished event: %#v", got)
	}
}

func TestSceneStepObserverEmitsLifecycleEvents(t *testing.T) {
	observer := &recordingObserver{}
	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Stay(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		observer:           observer,
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 23,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 10,
			Text:      "hello there",
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if len(observer.started) != 1 {
		t.Fatalf("expected one scene started event, got %d", len(observer.started))
	}
	if got := observer.started[0]; got.HandlerKind != HandlerSceneStepKind || got.HandlerName != "start" || got.Plugin != "wizard" {
		t.Fatalf("unexpected scene step started event: %#v", got)
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one scene finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; got.HandlerKind != HandlerSceneStepKind || got.HandlerName != "start" || got.Plugin != "wizard" || got.Err != nil {
		t.Fatalf("unexpected scene step finished event: %#v", got)
	}
}

func TestSceneMessageObserverEmitsLifecycleEvents(t *testing.T) {
	observer := &recordingObserver{}
	plugin := NewPlugin[NoData]("wizard")
	scene := plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Stay(), nil
		}).
		OnMessage(func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Exit(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		observer:           observer,
	}
	bot.AddPlugins(plugin)

	key, ok := buildSceneKey(SceneScopeUserChat, &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID: 42,
	})
	if !ok {
		t.Fatal("expected scene key to be built")
	}
	if err := bot.sessionStore.Set(key, SceneSession{Scene: scene.Name}); err != nil {
		t.Fatalf("failed to seed scene session: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 24,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 11,
			Text:      "hello there",
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if len(observer.started) != 1 {
		t.Fatalf("expected one scene started event, got %d", len(observer.started))
	}
	if got := observer.started[0]; got.HandlerKind != HandlerSceneMessageKind || got.HandlerName != "message_fallback" || got.Plugin != "wizard" {
		t.Fatalf("unexpected scene message started event: %#v", got)
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one scene finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; got.HandlerKind != HandlerSceneMessageKind || got.HandlerName != "message_fallback" || got.Plugin != "wizard" || got.Err != nil {
		t.Fatalf("unexpected scene message finished event: %#v", got)
	}
}

func TestScenePayloadHandlerRunsBeforeStep(t *testing.T) {
	payloadCalled := false
	stepCalled := false

	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			stepCalled = true
			return ctx.Stay(), nil
		}).
		OnPayload("confirm", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			payloadCalled = true
			if got, want := ctx.Args, []string{"7", "ok"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
				t.Fatalf("unexpected payload args: got %v want %v", got, want)
			}
			if ctx.Text != "" {
				t.Fatalf("callback flow must not populate Text, got %q", ctx.Text)
			}
			return ctx.Exit(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		payloadType:        BotPayloadJSON,
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	data, err := encodeJSONPayload(CallbackData{Command: "confirm", Args: []string{"7", "ok"}})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 25,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-scene",
			Data: data,
			From: tgapi.User{ID: 42},
			Message: &tgapi.Message{
				MessageID: 12,
				Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			},
		},
	})

	if !payloadCalled {
		t.Fatal("expected scene payload handler to be called")
	}
	if stepCalled {
		t.Fatal("expected scene payload to short-circuit the active step")
	}
}

func TestScenePayloadObserverEmitsLifecycleEvents(t *testing.T) {
	observer := &recordingObserver{}
	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Stay(), nil
		}).
		OnPayload("confirm", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Exit(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		payloadType:        BotPayloadJSON,
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		observer:           observer,
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	data, err := encodeJSONPayload(CallbackData{Command: "confirm"})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 26,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-scene",
			Data: data,
			From: tgapi.User{ID: 42},
			Message: &tgapi.Message{
				MessageID: 13,
				Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			},
		},
	})

	if len(observer.started) != 1 {
		t.Fatalf("expected one scene started event, got %d", len(observer.started))
	}
	if got := observer.started[0]; got.HandlerKind != HandlerScenePayloadKind || got.HandlerName != "confirm" || got.Plugin != "wizard" {
		t.Fatalf("unexpected scene payload started event: %#v", got)
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one scene finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; got.HandlerKind != HandlerScenePayloadKind || got.HandlerName != "confirm" || got.Plugin != "wizard" || got.Err != nil {
		t.Fatalf("unexpected scene payload finished event: %#v", got)
	}
}

func TestSceneUnmatchedPayloadFallsThroughWithoutRunningStep(t *testing.T) {
	stepCalled := false

	plugin := NewPlugin[NoData]("wizard")
	plugin.Payload("ping", func(ctx *MsgContext, db NoData) error { return nil })
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			stepCalled = true
			return ctx.Stay(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		payloadType:        BotPayloadJSON,
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	key, ok := buildSceneKey(SceneScopeUserChat, &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID: 42,
	})
	if !ok {
		t.Fatal("expected scene key to be built")
	}

	data, err := encodeJSONPayload(CallbackData{Command: "ping"})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 27,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-global",
			Data: data,
			From: tgapi.User{ID: 42},
			Message: &tgapi.Message{
				MessageID: 14,
				Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			},
		},
	})

	if stepCalled {
		t.Fatal("scene step must not run for an unmatched payload")
	}

	after, err := bot.sessionStore.Get(key)
	if err != nil {
		t.Fatalf("Get after handle returned error: %v", err)
	}
	if after.Scene != "signup" || after.Step != "start" {
		t.Fatalf("unexpected session after payload fallback: %#v", after)
	}
}

func TestScenePassDoesNotPersistSessionData(t *testing.T) {
	commandCalled := false

	plugin := NewPlugin[NoData]("wizard")
	plugin.Command("ping", func(ctx *MsgContext, db NoData) error {
		commandCalled = true
		return nil
	})
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			if err := ctx.SaveData(struct {
				Value string `json:"value"`
			}{Value: "changed"}); err != nil {
				t.Fatalf("SaveData returned error: %v", err)
			}
			return ctx.Pass(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	key, ok := buildSceneKey(SceneScopeUserChat, &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID: 42,
	})
	if !ok {
		t.Fatal("expected scene key to be built")
	}

	before, err := bot.sessionStore.Get(key)
	if err != nil {
		t.Fatalf("Get before handle returned error: %v", err)
	}
	if before.HasData() {
		t.Fatalf("expected empty session data before handle, got %#v", before)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 3,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 9,
			Text:      "/ping",
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if !commandCalled {
		t.Fatal("expected normal command routing to continue after SceneActionPass")
	}

	after, err := bot.sessionStore.Get(key)
	if err != nil {
		t.Fatalf("Get after handle returned error: %v", err)
	}
	if after.Scene != "signup" || after.Step != "start" {
		t.Fatalf("unexpected session after pass: %#v", after)
	}
	if after.HasData() {
		t.Fatalf("expected SceneActionPass to leave session data unchanged, got %#v", after)
	}
}

func TestSceneUnmatchedCommandFallsThroughWithoutRunningStep(t *testing.T) {
	commandCalled := false
	stepCalled := false

	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			stepCalled = true
			return ctx.Stay(), nil
		})
	plugin.Command("ping", func(ctx *MsgContext, db NoData) error {
		commandCalled = true
		return nil
	})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	key, ok := buildSceneKey(SceneScopeUserChat, &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID: 42,
	})
	if !ok {
		t.Fatal("expected scene key to be built")
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 5,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 10,
			Text:      "/ping",
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if !commandCalled {
		t.Fatal("expected normal command routing to handle /ping")
	}
	if stepCalled {
		t.Fatal("scene step must not run for an unmatched slash-command")
	}

	after, err := bot.sessionStore.Get(key)
	if err != nil {
		t.Fatalf("Get after handle returned error: %v", err)
	}
	if after.Scene != "signup" || after.Step != "start" {
		t.Fatalf("unexpected session after command fallback: %#v", after)
	}
}

func TestSceneMessageFallbackRunsWhenNoCommandOrStepMatch(t *testing.T) {
	fallbackCalled := false

	plugin := NewPlugin[NoData]("wizard")
	plugin.Scene("signup").
		SetEntry("start").
		OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Stay(), nil
		}).
		OnMessage(func(ctx *SceneContext, db NoData) (SceneResult, error) {
			fallbackCalled = true
			if ctx.Text != "hello fallback" {
				t.Fatalf("unexpected fallback text: got %q want %q", ctx.Text, "hello fallback")
			}
			return ctx.Exit(), nil
		})

	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		prefixes:           []string{"/"},
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
	}
	bot.AddPlugins(plugin)

	enterCtx := &MsgContext{
		Msg:          &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID:       42,
		sceneRuntime: bot,
	}
	if err := enterCtx.EnterScene("signup"); err != nil {
		t.Fatalf("EnterScene returned error: %v", err)
	}

	key, ok := buildSceneKey(SceneScopeUserChat, &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate}},
		FromID: 42,
	})
	if !ok {
		t.Fatal("expected scene key to be built")
	}
	if err := bot.sessionStore.Set(key, SceneSession{Scene: "signup", Step: "unknown"}); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 4,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 10,
			Text:      "hello fallback",
			Chat:      &tgapi.Chat{ID: 100, Type: tgapi.ChatTypePrivate},
			From:      &tgapi.User{ID: 42},
		},
	})

	if !fallbackCalled {
		t.Fatal("expected scene fallback handler to be called")
	}
}

func TestFindSceneSessionSupportsUserScopeWithoutMessage(t *testing.T) {
	bot := &Bot[NoData]{
		logger:             sneklog.NewLogger(),
		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUser, SceneScopeChat, SceneScopeUserChat},
	}

	if err := bot.sessionStore.Set("user_id:42", SceneSession{Scene: "signup", Step: "start"}); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	key, session, err := bot.findSceneSession(&MsgContext{FromID: 42})
	if err != nil {
		t.Fatalf("findSceneSession returned error: %v", err)
	}
	if key != "user_id:42" {
		t.Fatalf("unexpected session key: got %q want %q", key, "user_id:42")
	}
	if session.Scene != "signup" || session.Step != "start" {
		t.Fatalf("unexpected session: %#v", session)
	}
}

func TestSceneStoreErrorsPropagate(t *testing.T) {
	getErr := errors.New("get failed")
	setErr := errors.New("set failed")

	t.Run("find scene session get error", func(t *testing.T) {
		bot := &Bot[NoData]{
			logger:             sneklog.NewLogger(),
			sessionStore:       failingSessionStore{getErr: getErr},
			sceneScopePriority: []SceneScope{SceneScopeUser},
		}

		_, _, err := bot.findSceneSession(&MsgContext{FromID: 42})
		if !errors.Is(err, getErr) {
			t.Fatalf("expected getErr, got %v", err)
		}
	})

	t.Run("apply scene result set error", func(t *testing.T) {
		scene := NewScene[NoData]("signup").OnStep("start", func(ctx *SceneContext, db NoData) (SceneResult, error) {
			return ctx.Stay(), nil
		})
		bot := &Bot[NoData]{
			logger:             sneklog.NewLogger(),
			sessionStore:       failingSessionStore{setErr: setErr},
			sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
		}

		_, err := bot.applySceneResult(scene, &SceneContext{
			MsgContext: &MsgContext{},
			sess:       SceneSession{Scene: "signup", Step: "start"},
			key:        "user_id:42:chat_id:100",
		}, SceneResult{Action: SceneActionStay})
		if !errors.Is(err, setErr) {
			t.Fatalf("expected setErr, got %v", err)
		}
	})
}
