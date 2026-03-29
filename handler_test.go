package laniakea

import (
	"context"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
)

func ptr[T any](v T) *T {
	return &v
}

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

func TestPrepareUpdateCtxContract(t *testing.T) {
	tests := []struct {
		name               string
		update             *tgapi.Update
		wantMsg            bool
		wantFrom           bool
		wantFromID         int64
		wantCallbackID     string
		wantCallbackMsgID  int
		wantInlineMsgID    string
	}{
		{
			name: "message",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeMessage,
				Message: &tgapi.Message{
					MessageID: 11,
					From:      &tgapi.User{ID: 101},
					Chat:      &tgapi.Chat{ID: 1001},
				},
			},
			wantMsg:    true,
			wantFrom:   true,
			wantFromID: 101,
		},
		{
			name: "edited message",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeEditedMessage,
				EditedMessage: &tgapi.Message{
					MessageID: 12,
					From:      &tgapi.User{ID: 102},
					Chat:      &tgapi.Chat{ID: 1002},
				},
			},
			wantMsg:    true,
			wantFrom:   true,
			wantFromID: 102,
		},
		{
			name: "channel post sender chat",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeChannelPost,
				ChannelPost: &tgapi.Message{
					MessageID: 13,
					Chat:      &tgapi.Chat{ID: -1003},
				},
			},
			wantMsg: true,
		},
		{
			name: "business message",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeBusinessMessage,
				BusinessMessage: &tgapi.Message{
					MessageID: 14,
					From:      &tgapi.User{ID: 103},
					Chat:      &tgapi.Chat{ID: 1004},
				},
			},
			wantMsg:    true,
			wantFrom:   true,
			wantFromID: 103,
		},
		{
			name: "inline query",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeInlineQuery,
				InlineQuery: &tgapi.InlineQuery{ID: "iq", From: tgapi.User{ID: 104}},
			},
			wantFrom:   true,
			wantFromID: 104,
		},
		{
			name: "chosen inline result",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeChosenInlineResult,
				ChosenInlineResult: &tgapi.ChosenInlineResult{ResultID: "res", From: tgapi.User{ID: 105}},
			},
			wantFrom:   true,
			wantFromID: 105,
		},
		{
			name: "callback query with message",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeCallbackQuery,
				CallbackQuery: &tgapi.CallbackQuery{
					ID:   "cb-1",
					From: tgapi.User{ID: 106},
					Message: &tgapi.Message{
						MessageID: 77,
						Chat:      &tgapi.Chat{ID: 1005},
					},
				},
			},
			wantMsg:           true,
			wantFrom:          true,
			wantFromID:        106,
			wantCallbackID:    "cb-1",
			wantCallbackMsgID: 77,
		},
		{
			name: "callback query with inline message",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeCallbackQuery,
				CallbackQuery: &tgapi.CallbackQuery{
					ID:              "cb-2",
					From:            tgapi.User{ID: 107},
					InlineMessageID: ptr("inline-42"),
				},
			},
			wantFrom:       true,
			wantFromID:     107,
			wantCallbackID: "cb-2",
			wantInlineMsgID:"inline-42",
		},
		{
			name: "shipping query",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeShippingQuery,
				ShippingQuery: &tgapi.ShippingQuery{ID: "ship", From: tgapi.User{ID: 108}},
			},
			wantFrom:   true,
			wantFromID: 108,
		},
		{
			name: "pre checkout query",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypePreCheckoutQuery,
				PreCheckoutQuery: &tgapi.PreCheckoutQuery{ID: "pre", From: tgapi.User{ID: 109}},
			},
			wantFrom:   true,
			wantFromID: 109,
		},
		{
			name: "purchased paid media",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypePurchasedPaidMedia,
				PurchasedPaidMedia: &tgapi.PaidMediaPurchased{From: tgapi.User{ID: 110}},
			},
			wantFrom:   true,
			wantFromID: 110,
		},
		{
			name: "my chat member",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeMyChatMember,
				MyChatMember: &tgapi.ChatMemberUpdated{From: tgapi.User{ID: 111}},
			},
			wantFrom:   true,
			wantFromID: 111,
		},
		{
			name: "chat member",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeChatMember,
				ChatMember: &tgapi.ChatMemberUpdated{From: tgapi.User{ID: 112}},
			},
			wantFrom:   true,
			wantFromID: 112,
		},
		{
			name: "chat join request",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeChatJoinRequest,
				ChatJoinRequest: &tgapi.ChatJoinRequest{From: tgapi.User{ID: 113}},
			},
			wantFrom:   true,
			wantFromID: 113,
		},
		{
			name: "business connection",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeBusinessConnection,
				BusinessConnection: &tgapi.BusinessConnection{User: tgapi.User{ID: 114}},
			},
			wantFrom:   true,
			wantFromID: 114,
		},
		{
			name: "poll answer",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypePollAnswer,
				PollAnswer: &tgapi.PollAnswer{User: tgapi.User{ID: 115}},
			},
			wantFrom:   true,
			wantFromID: 115,
		},
		{
			name: "message reaction",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeMessageReaction,
				MessageReaction: &tgapi.MessageReactionUpdated{User: &tgapi.User{ID: 116}},
			},
			wantFrom:   true,
			wantFromID: 116,
		},
		{
			name: "chat boost",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeChatBoost,
				ChatBoost: &tgapi.ChatBoostUpdated{
					Boost: tgapi.ChatBoost{Source: tgapi.ChatBoostSource{User: tgapi.User{ID: 117}}},
				},
			},
			wantFrom:   true,
			wantFromID: 117,
		},
		{
			name: "removed chat boost",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeRemovedChatBoost,
				RemovedChatBoost: &tgapi.ChatBoostRemoved{
					Source: tgapi.ChatBoostSource{User: tgapi.User{ID: 118}},
				},
			},
			wantFrom:   true,
			wantFromID: 118,
		},
		{
			name: "poll",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypePoll,
				Poll: &tgapi.Poll{ID: "poll"},
			},
		},
		{
			name: "message reaction count",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeMessageReactionCount,
				MessageReactionCount: &tgapi.MessageReactionCountUpdated{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot[NoDB]{}
			ctx := &MsgContext{}
			bot.prepareUpdateCtx(tt.update, ctx)

			if got := ctx.Msg != nil; got != tt.wantMsg {
				t.Fatalf("unexpected Msg presence: got %v want %v", got, tt.wantMsg)
			}
			if got := ctx.From != nil; got != tt.wantFrom {
				t.Fatalf("unexpected From presence: got %v want %v", got, tt.wantFrom)
			}
			if ctx.FromID != tt.wantFromID {
				t.Fatalf("unexpected FromID: got %d want %d", ctx.FromID, tt.wantFromID)
			}
			if ctx.CallbackQueryId != tt.wantCallbackID {
				t.Fatalf("unexpected CallbackQueryId: got %q want %q", ctx.CallbackQueryId, tt.wantCallbackID)
			}
			if ctx.CallbackMsgId != tt.wantCallbackMsgID {
				t.Fatalf("unexpected CallbackMsgId: got %d want %d", ctx.CallbackMsgId, tt.wantCallbackMsgID)
			}
			if ctx.InlineMsgId != tt.wantInlineMsgID {
				t.Fatalf("unexpected InlineMsgId: got %q want %q", ctx.InlineMsgId, tt.wantInlineMsgID)
			}
			if ctx.Text != "" {
				t.Fatalf("prepareUpdateCtx must not populate Text, got %q", ctx.Text)
			}
			if len(ctx.Args) != 0 {
				t.Fatalf("prepareUpdateCtx must not populate Args, got %v", ctx.Args)
			}
		})
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

func TestHandleEditedMessageStaysOutOfCommandFlow(t *testing.T) {
	commandCalled := false
	updateCalled := false

	plugin := NewPlugin[NoDB]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoDB) error {
		commandCalled = true
		return nil
	}, "ping")
	plugin.AddUpdateHandler(tgapi.UpdateTypeEditedMessage, func(ctx *MsgContext, db NoDB) error {
		updateCalled = true
		if ctx.Msg == nil {
			t.Fatal("expected ctx.Msg in edited message handler")
		}
		if ctx.Text != "" {
			t.Fatalf("expected empty Text in edited_message update handler, got %q", ctx.Text)
		}
		if len(ctx.Args) != 0 {
			t.Fatalf("expected empty Args in edited_message update handler, got %v", ctx.Args)
		}
		return nil
	})

	bot := &Bot[NoDB]{
		logger:   slog.CreateLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoDB]{clonePlugin(plugin)},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 20,
		Type:     tgapi.UpdateTypeEditedMessage,
		EditedMessage: &tgapi.Message{
			MessageID: 1,
			Text:      "/ping",
			From:      &tgapi.User{ID: 1},
			Chat:      &tgapi.Chat{ID: 42},
		},
	})

	if commandCalled {
		t.Fatal("edited_message must not enter command flow")
	}
	if !updateCalled {
		t.Fatal("expected edited_message update handler to be called")
	}
}

func TestHandleEditedChannelPostStaysOutOfCommandFlow(t *testing.T) {
	commandCalled := false
	updateCalled := false

	plugin := NewPlugin[NoDB]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoDB) error {
		commandCalled = true
		return nil
	}, "ping")
	plugin.AddUpdateHandler(tgapi.UpdateTypeEditedChannelPost, func(ctx *MsgContext, db NoDB) error {
		updateCalled = true
		if ctx.Msg == nil {
			t.Fatal("expected ctx.Msg in edited channel post handler")
		}
		return nil
	})

	bot := &Bot[NoDB]{
		logger:   slog.CreateLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoDB]{clonePlugin(plugin)},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 21,
		Type:     tgapi.UpdateTypeEditedChannelPost,
		EditedChannelPost: &tgapi.Message{
			MessageID: 1,
			Text:      "/ping",
			Chat:      &tgapi.Chat{ID: -10042},
		},
	})

	if commandCalled {
		t.Fatal("edited_channel_post must not enter command flow")
	}
	if !updateCalled {
		t.Fatal("expected edited_channel_post update handler to be called")
	}
}

func TestHandleCallbackPopulatesMessageTargets(t *testing.T) {
	called := false
	plugin := NewPlugin[NoDB]("test")
	plugin.NewPayload(func(ctx *MsgContext, db NoDB) error {
		called = true
		if ctx.CallbackQueryId != "cb-msg" {
			t.Fatalf("unexpected CallbackQueryId: %q", ctx.CallbackQueryId)
		}
		if ctx.CallbackMsgId != 55 {
			t.Fatalf("unexpected CallbackMsgId: %d", ctx.CallbackMsgId)
		}
		if ctx.InlineMsgId != "" {
			t.Fatalf("did not expect InlineMsgId, got %q", ctx.InlineMsgId)
		}
		if ctx.Msg == nil {
			t.Fatal("expected callback message context")
		}
		if ctx.From == nil || ctx.FromID != 7 {
			t.Fatalf("unexpected callback sender: %#v / %d", ctx.From, ctx.FromID)
		}
		if ctx.Text != "" {
			t.Fatalf("callback flow must not populate Text, got %q", ctx.Text)
		}
		if got, want := ctx.Args, []string{"7", "ok"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("unexpected callback args: got %v want %v", got, want)
		}
		return nil
	}, "approve")

	bot := &Bot[NoDB]{
		logger:      slog.CreateLogger(),
		payloadType: BotPayloadJson,
		plugins:     []Plugin[NoDB]{clonePlugin(plugin)},
	}

	data, err := encodeJsonPayload(CallbackData{Command: "approve", Args: []string{"7", "ok"}})
	if err != nil {
		t.Fatalf("encodeJsonPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 30,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-msg",
			Data: data,
			From: tgapi.User{ID: 7},
			Message: &tgapi.Message{
				MessageID: 55,
				Chat:      &tgapi.Chat{ID: 77},
			},
		},
	})

	if !called {
		t.Fatal("expected payload handler to be called")
	}
}

func TestHandleCallbackPopulatesInlineTargets(t *testing.T) {
	called := false
	plugin := NewPlugin[NoDB]("test")
	plugin.NewPayload(func(ctx *MsgContext, db NoDB) error {
		called = true
		if ctx.CallbackQueryId != "cb-inline" {
			t.Fatalf("unexpected CallbackQueryId: %q", ctx.CallbackQueryId)
		}
		if ctx.CallbackMsgId != 0 {
			t.Fatalf("did not expect CallbackMsgId, got %d", ctx.CallbackMsgId)
		}
		if ctx.InlineMsgId != "inline-55" {
			t.Fatalf("unexpected InlineMsgId: %q", ctx.InlineMsgId)
		}
		if ctx.Msg != nil {
			t.Fatalf("did not expect callback chat message context, got %#v", ctx.Msg)
		}
		if ctx.From == nil || ctx.FromID != 8 {
			t.Fatalf("unexpected callback sender: %#v / %d", ctx.From, ctx.FromID)
		}
		if ctx.Text != "" {
			t.Fatalf("callback flow must not populate Text, got %q", ctx.Text)
		}
		if got, want := ctx.Args, []string{"9"}; len(got) != len(want) || got[0] != want[0] {
			t.Fatalf("unexpected callback args: got %v want %v", got, want)
		}
		return nil
	}, "inline.approve")

	bot := &Bot[NoDB]{
		logger:      slog.CreateLogger(),
		payloadType: BotPayloadJson,
		plugins:     []Plugin[NoDB]{clonePlugin(plugin)},
	}

	data, err := encodeJsonPayload(CallbackData{Command: "inline.approve", Args: []string{"9"}})
	if err != nil {
		t.Fatalf("encodeJsonPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 31,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:              "cb-inline",
			Data:            data,
			From:            tgapi.User{ID: 8},
			InlineMessageID: ptr("inline-55"),
		},
	})

	if !called {
		t.Fatal("expected inline payload handler to be called")
	}
}
