package laniakea

import (
	"context"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
)

func TestCheckPrefixesSkipsEmptyPrefixes(t *testing.T) {
	bot := &Bot[NoDB]{prefixes: []string{"", "/"}}

	if prefix, ok := bot.checkPrefixes("hello"); ok {
		t.Fatalf("unexpected prefix match for plain text: %q", prefix)
	}
	if prefix, ok := bot.checkPrefixes("/start"); !ok || prefix != "/" {
		t.Fatalf("unexpected prefix result: prefix=%q ok=%v", prefix, ok)
	}
}

func TestBotMiddlewareReceivesLogger(t *testing.T) {
	logger := slog.CreateLogger()
	called := false

	bot := &Bot[NoDB]{
		logger: logger,
		middlewares: []Middleware[NoDB]{
			NewMiddleware("logger-check", func(ctx *MsgContext, db NoDB) bool {
				called = true
				if ctx.Logger != logger {
					t.Fatalf("expected bot logger in middleware context, got %#v", ctx.Logger)
				}
				return true
			}),
		},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 1,
		Type:     tgapi.UpdateTypePoll,
		Poll: &tgapi.Poll{
			ID:       "poll",
			Question: "question",
		},
	})

	if !called {
		t.Fatal("expected bot middleware to be called")
	}
}

func TestAddUpdateHandlerRejectsReservedUpdateTypes(t *testing.T) {
	plugin := NewPlugin[NoDB]("test")
	handler := func(ctx *MsgContext, db NoDB) error { return nil }

	for _, updateType := range []tgapi.UpdateType{
		tgapi.UpdateTypeMessage,
		tgapi.UpdateTypeChannelPost,
		tgapi.UpdateTypeCallbackQuery,
	} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("AddUpdateHandler(%q) panicked: %v", updateType, r)
				}
			}()
			plugin.AddUpdateHandler(updateType, handler)
		}()

		if _, ok := plugin.handlers[updateType]; ok {
			t.Fatalf("reserved update type %q must not be registered", updateType)
		}
	}
}

func TestHandleUpdateHandlersPopulateFromContext(t *testing.T) {
	tests := []struct {
		name   string
		update *tgapi.Update
		wantID int64
	}{
		{
			name: "inline query",
			update: &tgapi.Update{
				UpdateID: 1,
				Type:     tgapi.UpdateTypeInlineQuery,
				InlineQuery: &tgapi.InlineQuery{
					ID:    "iq",
					From:  tgapi.User{ID: 41},
					Query: "ping",
				},
			},
			wantID: 41,
		},
		{
			name: "chosen inline result",
			update: &tgapi.Update{
				UpdateID: 2,
				Type:     tgapi.UpdateTypeChosenInlineResult,
				ChosenInlineResult: &tgapi.ChosenInlineResult{
					ResultID: "res",
					From:     tgapi.User{ID: 77},
					Query:    "pong",
				},
			},
			wantID: 77,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			plugin := NewPlugin[NoDB]("test").AddUpdateHandler(tt.update.Type, func(ctx *MsgContext, db NoDB) error {
				called = true
				if ctx.Update.UpdateID != tt.update.UpdateID {
					t.Fatalf("unexpected update in context: got %d want %d", ctx.Update.UpdateID, tt.update.UpdateID)
				}
				if ctx.From == nil {
					t.Fatal("expected ctx.From to be populated")
				}
				if ctx.FromID != tt.wantID {
					t.Fatalf("unexpected FromID: got %d want %d", ctx.FromID, tt.wantID)
				}
				if ctx.From.ID != tt.wantID {
					t.Fatalf("unexpected ctx.From.ID: got %d want %d", ctx.From.ID, tt.wantID)
				}
				if ctx.Msg != nil {
					t.Fatalf("did not expect message context for %s", tt.name)
				}
				return nil
			})

			bot := &Bot[NoDB]{
				logger:  slog.CreateLogger(),
				plugins: []Plugin[NoDB]{clonePlugin(plugin)},
			}

			bot.handle(context.Background(), tt.update)

			if !called {
				t.Fatalf("expected update handler for %s to be called", tt.name)
			}
		})
	}
}

func TestHandleUpdateHandlersReceiveIsolatedContexts(t *testing.T) {
	firstCalled := false
	secondCalled := false

	first := NewPlugin[NoDB]("first").AddUpdateHandler(tgapi.UpdateTypeInlineQuery, func(ctx *MsgContext, db NoDB) error {
		firstCalled = true
		if ctx.FromID != 41 {
			t.Fatalf("unexpected FromID in first handler: got %d want 41", ctx.FromID)
		}
		ctx.From = nil
		ctx.FromID = 999
		ctx.Text = "mutated"
		ctx.Args = []string{"mutated"}
		return nil
	})
	second := NewPlugin[NoDB]("second").AddUpdateHandler(tgapi.UpdateTypeInlineQuery, func(ctx *MsgContext, db NoDB) error {
		secondCalled = true
		if ctx.From == nil {
			t.Fatal("expected ctx.From to remain populated for second handler")
		}
		if ctx.FromID != 41 {
			t.Fatalf("unexpected FromID in second handler: got %d want 41", ctx.FromID)
		}
		if ctx.Text != "" {
			t.Fatalf("unexpected leaked Text in second handler: %q", ctx.Text)
		}
		if len(ctx.Args) != 0 {
			t.Fatalf("unexpected leaked Args in second handler: %v", ctx.Args)
		}
		return nil
	})

	bot := &Bot[NoDB]{
		logger: slog.CreateLogger(),
		plugins: []Plugin[NoDB]{
			clonePlugin(first),
			clonePlugin(second),
		},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 3,
		Type:     tgapi.UpdateTypeInlineQuery,
		InlineQuery: &tgapi.InlineQuery{
			ID:    "iq",
			From:  tgapi.User{ID: 41},
			Query: "ping",
		},
	})

	if !firstCalled || !secondCalled {
		t.Fatalf("expected both handlers to be called, got first=%v second=%v", firstCalled, secondCalled)
	}
}

func TestHandleChannelPostCommandWithSenderChat(t *testing.T) {
	called := false
	plugin := NewPlugin[NoDB]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoDB) error {
		called = true
		if ctx.Msg == nil {
			t.Fatal("expected message context")
		}
		if ctx.Msg.Chat == nil || ctx.Msg.Chat.ID != -1001 {
			t.Fatalf("unexpected chat context: %#v", ctx.Msg.Chat)
		}
		if ctx.From != nil {
			t.Fatalf("expected ctx.From to stay nil for sender_chat updates, got %#v", ctx.From)
		}
		if ctx.FromID != 0 {
			t.Fatalf("expected zero FromID for sender_chat updates, got %d", ctx.FromID)
		}
		return nil
	}, "ping")

	bot := &Bot[NoDB]{
		logger:   slog.CreateLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoDB]{clonePlugin(plugin)},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 10,
		Type:     tgapi.UpdateTypeChannelPost,
		ChannelPost: &tgapi.Message{
			MessageID:  55,
			Text:       "/ping",
			SenderChat: &tgapi.Chat{ID: -1001, Type: string(tgapi.ChatTypeChannel)},
			Chat:       &tgapi.Chat{ID: -1001, Type: string(tgapi.ChatTypeChannel)},
		},
	})

	if !called {
		t.Fatal("expected channel post command handler to be called")
	}
}

func TestCommandHandlerBindArgsEndToEnd(t *testing.T) {
	type banInput struct {
		UserID int
		Reason string
	}

	var got banInput
	plugin := NewPlugin[NoDB]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoDB) error {
		return ctx.BindArgs(&got)
	}, "ban",
		NewCommandArg("user_id").SetValueType(CommandValueIntType).SetRequired(),
		NewCommandArg("reason").SetRequired(),
	)

	bot := &Bot[NoDB]{
		logger:   slog.CreateLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoDB]{clonePlugin(plugin)},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 11,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 1,
			Text:      "/ban 42 too loud",
			Chat:      &tgapi.Chat{ID: 99, Type: string(tgapi.ChatTypePrivate)},
		},
	})

	want := banInput{UserID: 42, Reason: "too loud"}
	if got != want {
		t.Fatalf("unexpected bound input: got %#v want %#v", got, want)
	}
}

func TestPayloadHandlerBindArgsEndToEnd(t *testing.T) {
	type payloadInput struct {
		ID   int
		Note string
	}

	var got payloadInput
	plugin := NewPlugin[NoDB]("test")
	plugin.NewPayload(func(ctx *MsgContext, db NoDB) error {
		return ctx.BindArgs(&got)
	}, "approve",
		NewCommandArg("id").SetValueType(CommandValueIntType).SetRequired(),
		NewCommandArg("note").SetRequired(),
	)

	bot := &Bot[NoDB]{
		logger:      slog.CreateLogger(),
		payloadType: BotPayloadJson,
		plugins:     []Plugin[NoDB]{clonePlugin(plugin)},
	}

	data, err := encodeJsonPayload(CallbackData{
		Command: "approve",
		Args:    []string{"7", "looks", "good"},
	})
	if err != nil {
		t.Fatalf("encodeJsonPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 12,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-1",
			Data: data,
			From: tgapi.User{ID: 1},
		},
	})

	want := payloadInput{ID: 7, Note: "looks good"}
	if got != want {
		t.Fatalf("unexpected bound payload input: got %#v want %#v", got, want)
	}
}
