package laniakea

import (
	"context"
	"errors"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

type recordingObserver struct {
	received []UpdateReceivedEvent
	started  []HandlerStartedEvent
	finished []HandlerFinishedEvent
	errors   []ErrorEvent
	handled  []UpdateHandledEvent
	policies []PolicyCheckedEvent
	runners  []RunnerFinishedEvent
	retries  []PollingRetryEvent
}

func (o *recordingObserver) OnReceiveUpdate(_ context.Context, ev UpdateReceivedEvent) {
	o.received = append(o.received, ev)
}
func (o *recordingObserver) OnHandledUpdate(_ context.Context, ev UpdateHandledEvent) {
	o.handled = append(o.handled, ev)
}
func (o *recordingObserver) OnHandlerStarted(_ context.Context, ev HandlerStartedEvent) {
	o.started = append(o.started, ev)
}
func (o *recordingObserver) OnHandlerFinished(_ context.Context, ev HandlerFinishedEvent) {
	o.finished = append(o.finished, ev)
}
func (*recordingObserver) OnSceneTransition(context.Context, SceneTransitionEvent) {}
func (o *recordingObserver) OnPolicyChecked(_ context.Context, ev PolicyCheckedEvent) {
	o.policies = append(o.policies, ev)
}
func (o *recordingObserver) OnRunnerFinished(_ context.Context, ev RunnerFinishedEvent) {
	o.runners = append(o.runners, ev)
}
func (o *recordingObserver) OnPollingRetry(_ context.Context, ev PollingRetryEvent) {
	o.retries = append(o.retries, ev)
}
func (o *recordingObserver) OnError(_ context.Context, ev ErrorEvent) {
	o.errors = append(o.errors, ev)
}

func TestCheckPrefixesSkipsEmptyPrefixes(t *testing.T) {
	bot := &Bot[NoData]{prefixes: []string{"", "/"}}

	if prefix, ok := bot.checkPrefixes("hello"); ok {
		t.Fatalf("unexpected prefix match for plain text: %q", prefix)
	}
	if prefix, ok := bot.checkPrefixes("/start"); !ok || prefix != "/" {
		t.Fatalf("unexpected prefix result: prefix=%q ok=%v", prefix, ok)
	}
}

func TestBotMiddlewareReceivesLogger(t *testing.T) {
	logger := sneklog.NewLogger()
	called := false

	bot := &Bot[NoData]{
		logger: logger,
		middlewares: []Middleware[NoData]{
			NewMiddleware("logger-check", func(ctx *MsgContext, db NoData) bool {
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
	plugin := NewPlugin[NoData]("test")
	handler := func(ctx *MsgContext, db NoData) error { return nil }

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
		name              string
		update            *tgapi.Update
		wantMsg           bool
		wantFrom          bool
		wantFromID        int64
		wantChat          bool
		wantChatID        int64
		wantCallbackID    string
		wantCallbackMsgID int
		wantInlineMsgID   string
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
			wantChat:   true,
			wantChatID: 1001,
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
			wantChat:   true,
			wantChatID: 1002,
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
			wantMsg:    true,
			wantChat:   true,
			wantChatID: -1003,
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
			wantChat:   true,
			wantChatID: 1004,
		},
		{
			name: "inline query",
			update: &tgapi.Update{
				Type:        tgapi.UpdateTypeInlineQuery,
				InlineQuery: &tgapi.InlineQuery{ID: "iq", From: tgapi.User{ID: 104}},
			},
			wantFrom:   true,
			wantFromID: 104,
		},
		{
			name: "chosen inline result",
			update: &tgapi.Update{
				Type:               tgapi.UpdateTypeChosenInlineResult,
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
			wantChat:          true,
			wantChatID:        1005,
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
					InlineMessageID: new("inline-42"),
				},
			},
			wantFrom:        true,
			wantFromID:      107,
			wantCallbackID:  "cb-2",
			wantInlineMsgID: "inline-42",
		},
		{
			name: "shipping query",
			update: &tgapi.Update{
				Type:          tgapi.UpdateTypeShippingQuery,
				ShippingQuery: &tgapi.ShippingQuery{ID: "ship", From: tgapi.User{ID: 108}},
			},
			wantFrom:   true,
			wantFromID: 108,
		},
		{
			name: "pre checkout query",
			update: &tgapi.Update{
				Type:             tgapi.UpdateTypePreCheckoutQuery,
				PreCheckoutQuery: &tgapi.PreCheckoutQuery{ID: "pre", From: tgapi.User{ID: 109}},
			},
			wantFrom:   true,
			wantFromID: 109,
		},
		{
			name: "purchased paid media",
			update: &tgapi.Update{
				Type:               tgapi.UpdateTypePurchasedPaidMedia,
				PurchasedPaidMedia: &tgapi.PaidMediaPurchased{From: tgapi.User{ID: 110}},
			},
			wantFrom:   true,
			wantFromID: 110,
		},
		{
			name: "my chat member",
			update: &tgapi.Update{
				Type:         tgapi.UpdateTypeMyChatMember,
				MyChatMember: &tgapi.ChatMemberUpdated{From: tgapi.User{ID: 111}, Chat: tgapi.Chat{ID: -2001}},
			},
			wantFrom:   true,
			wantFromID: 111,
			wantChat:   true,
			wantChatID: -2001,
		},
		{
			name: "chat member",
			update: &tgapi.Update{
				Type:       tgapi.UpdateTypeChatMember,
				ChatMember: &tgapi.ChatMemberUpdated{From: tgapi.User{ID: 112}, Chat: tgapi.Chat{ID: -2002}},
			},
			wantFrom:   true,
			wantFromID: 112,
			wantChat:   true,
			wantChatID: -2002,
		},
		{
			name: "chat join request",
			update: &tgapi.Update{
				Type:            tgapi.UpdateTypeChatJoinRequest,
				ChatJoinRequest: &tgapi.ChatJoinRequest{From: tgapi.User{ID: 113}, Chat: tgapi.Chat{ID: -2003}},
			},
			wantFrom:   true,
			wantFromID: 113,
			wantChat:   true,
			wantChatID: -2003,
		},
		{
			name: "business connection",
			update: &tgapi.Update{
				Type:               tgapi.UpdateTypeBusinessConnection,
				BusinessConnection: &tgapi.BusinessConnection{User: tgapi.User{ID: 114}},
			},
			wantFrom:   true,
			wantFromID: 114,
		},
		{
			name: "poll answer",
			update: &tgapi.Update{
				Type:       tgapi.UpdateTypePollAnswer,
				PollAnswer: &tgapi.PollAnswer{User: tgapi.User{ID: 115}},
			},
			wantFrom:   true,
			wantFromID: 115,
		},
		{
			name: "message reaction",
			update: &tgapi.Update{
				Type:            tgapi.UpdateTypeMessageReaction,
				MessageReaction: &tgapi.MessageReactionUpdated{User: &tgapi.User{ID: 116}, Chat: &tgapi.Chat{ID: -2004}},
			},
			wantFrom:   true,
			wantFromID: 116,
			wantChat:   true,
			wantChatID: -2004,
		},
		{
			name: "chat boost",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeChatBoost,
				ChatBoost: &tgapi.ChatBoostUpdated{
					Chat:  tgapi.Chat{ID: -2005},
					Boost: tgapi.ChatBoost{Source: tgapi.ChatBoostSource{User: tgapi.User{ID: 117}}},
				},
			},
			wantFrom:   true,
			wantFromID: 117,
			wantChat:   true,
			wantChatID: -2005,
		},
		{
			name: "removed chat boost",
			update: &tgapi.Update{
				Type: tgapi.UpdateTypeRemovedChatBoost,
				RemovedChatBoost: &tgapi.ChatBoostRemoved{
					Chat:   tgapi.Chat{ID: -2006},
					Source: tgapi.ChatBoostSource{User: tgapi.User{ID: 118}},
				},
			},
			wantFrom:   true,
			wantFromID: 118,
			wantChat:   true,
			wantChatID: -2006,
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
				Type:                 tgapi.UpdateTypeMessageReactionCount,
				MessageReactionCount: &tgapi.MessageReactionCountUpdated{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot[NoData]{}
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
			if got := ctx.Chat != nil; got != tt.wantChat {
				t.Fatalf("unexpected Chat presence: got %v want %v", got, tt.wantChat)
			}
			if ctx.ChatID != tt.wantChatID {
				t.Fatalf("unexpected ChatID: got %d want %d", ctx.ChatID, tt.wantChatID)
			}
			if ctx.CallbackQueryID != tt.wantCallbackID {
				t.Fatalf("unexpected CallbackQueryID: got %q want %q", ctx.CallbackQueryID, tt.wantCallbackID)
			}
			if ctx.CallbackMsgID != tt.wantCallbackMsgID {
				t.Fatalf("unexpected CallbackMsgID: got %d want %d", ctx.CallbackMsgID, tt.wantCallbackMsgID)
			}
			if ctx.InlineMsgID != tt.wantInlineMsgID {
				t.Fatalf("unexpected InlineMsgID: got %q want %q", ctx.InlineMsgID, tt.wantInlineMsgID)
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
			plugin := NewPlugin[NoData]("test").AddUpdateHandler(tt.update.Type, func(ctx *MsgContext, db NoData) error {
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

			bot := &Bot[NoData]{
				logger:  sneklog.NewLogger(),
				plugins: []Plugin[NoData]{clonePlugin(plugin)},
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

	first := NewPlugin[NoData]("first").AddUpdateHandler(tgapi.UpdateTypeInlineQuery, func(ctx *MsgContext, db NoData) error {
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
	second := NewPlugin[NoData]("second").AddUpdateHandler(tgapi.UpdateTypeInlineQuery, func(ctx *MsgContext, db NoData) error {
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

	bot := &Bot[NoData]{
		logger: sneklog.NewLogger(),
		plugins: []Plugin[NoData]{
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

func TestHandleUpdateObserverEmitsUpdateErrors(t *testing.T) {
	observer := &recordingObserver{}
	plugin := NewPlugin[NoData]("test").AddUpdateHandler(tgapi.UpdateTypeInlineQuery, func(ctx *MsgContext, db NoData) error {
		return AsUserError(errors.New("update failed"))
	})

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		plugins:  []Plugin[NoData]{clonePlugin(plugin)},
		observer: observer,
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 4,
		Type:     tgapi.UpdateTypeInlineQuery,
		InlineQuery: &tgapi.InlineQuery{
			ID:   "iq",
			From: tgapi.User{ID: 41},
		},
	})

	if len(observer.errors) != 1 {
		t.Fatalf("expected one observer error event, got %d", len(observer.errors))
	}
	ev := observer.errors[0]
	if ev.Plugin != "test" {
		t.Fatalf("unexpected plugin: %q", ev.Plugin)
	}
	if ev.HandlerKind != HandlerUpdateKind {
		t.Fatalf("unexpected handler kind: %q", ev.HandlerKind)
	}
	if ev.HandlerName != string(tgapi.UpdateTypeInlineQuery) {
		t.Fatalf("unexpected handler name: %q", ev.HandlerName)
	}
	if !ev.UserFacing {
		t.Fatal("expected update error to be marked user-facing")
	}
	if len(observer.started) != 1 {
		t.Fatalf("expected one handler started event, got %d", len(observer.started))
	}
	if got := observer.started[0]; got.HandlerKind != HandlerUpdateKind || got.HandlerName != string(tgapi.UpdateTypeInlineQuery) || got.Plugin != "test" {
		t.Fatalf("unexpected started event: %#v", got)
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one handler finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; got.HandlerKind != HandlerUpdateKind || got.HandlerName != string(tgapi.UpdateTypeInlineQuery) || got.Plugin != "test" || got.Err == nil || !got.UserFacing {
		t.Fatalf("unexpected finished event: %#v", got)
	}
}

func TestHandleObserverCompletesUpdateWhenBotMiddlewareBlocks(t *testing.T) {
	observer := &recordingObserver{}
	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		observer: observer,
		middlewares: []Middleware[NoData]{
			NewMiddleware("block", func(ctx *MsgContext, db NoData) bool {
				return false
			}),
		},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 8,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 1,
			Date:      1,
			Chat:      &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate},
			Text:      "/start",
		},
	})

	if len(observer.received) != 1 {
		t.Fatalf("expected one received event, got %d", len(observer.received))
	}
	if len(observer.handled) != 1 {
		t.Fatalf("expected one handled event, got %d", len(observer.handled))
	}
	got := observer.handled[0]
	if got.UpdateID != 8 || got.UpdateType != tgapi.UpdateTypeMessage || got.ChatID != 42 || got.Handled {
		t.Fatalf("unexpected handled event: %#v", got)
	}
	if len(observer.started) != 0 || len(observer.finished) != 0 || len(observer.errors) != 0 {
		t.Fatalf("middleware block should not emit handler lifecycle or errors: started=%d finished=%d errors=%d", len(observer.started), len(observer.finished), len(observer.errors))
	}
}

func TestHandleMessageFallbackRunsAfterCommandMiss(t *testing.T) {
	observer := &recordingObserver{}
	called := false
	plugin := NewPlugin[NoData]("test")
	plugin.SetMessageFallback(func(ctx *MsgContext, db NoData) error {
		called = true
		if ctx.Text != "/missing hello world" {
			t.Fatalf("unexpected fallback text: got %q", ctx.Text)
		}
		if ctx.Prefix != "/" {
			t.Fatalf("unexpected fallback prefix: got %q", ctx.Prefix)
		}
		wantArgs := []string{"/missing", "hello", "world"}
		if len(ctx.Args) != len(wantArgs) || ctx.Args[0] != wantArgs[0] || ctx.Args[1] != wantArgs[1] || ctx.Args[2] != wantArgs[2] {
			t.Fatalf("unexpected fallback args: got %v want %v", ctx.Args, wantArgs)
		}
		return nil
	})

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
		observer: observer,
	}
	bot.AddPlugins(plugin)

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 5,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 1,
			Text:      "/missing hello world",
			From:      &tgapi.User{ID: 41},
			Chat:      &tgapi.Chat{ID: 99},
		},
	})

	if !called {
		t.Fatal("expected message fallback to be called")
	}
	if len(observer.started) != 1 {
		t.Fatalf("expected one started event, got %d", len(observer.started))
	}
	if got := observer.started[0]; got.HandlerKind != HandlerMessageKind || got.HandlerName != "message_fallback" || got.Plugin != "test" {
		t.Fatalf("unexpected started event: %#v", got)
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; got.HandlerKind != HandlerMessageKind || got.HandlerName != "message_fallback" || got.Plugin != "test" || got.Err != nil {
		t.Fatalf("unexpected finished event: %#v", got)
	}
	if len(observer.handled) != 1 || !observer.handled[0].Handled {
		t.Fatalf("expected handled update event, got %#v", observer.handled)
	}
}

func TestHandleMessageFallbackRunsForPlainText(t *testing.T) {
	called := false
	plugin := NewPlugin[NoData]("test").SetMessageFallback(func(ctx *MsgContext, db NoData) error {
		called = true
		if ctx.Text != "hello fallback" {
			t.Fatalf("unexpected fallback text: got %q", ctx.Text)
		}
		if ctx.Prefix != "" {
			t.Fatalf("unexpected fallback prefix: got %q", ctx.Prefix)
		}
		return nil
	})

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
	}
	bot.AddPlugins(plugin)

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 6,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 1,
			Text:      "hello fallback",
			From:      &tgapi.User{ID: 41},
			Chat:      &tgapi.Chat{ID: 99},
		},
	})

	if !called {
		t.Fatal("expected message fallback to be called")
	}
}

func TestHandleMessageFallbackRespectsMiddleware(t *testing.T) {
	called := false
	plugin := NewPlugin[NoData]("test")
	plugin.AddMiddleware(NewMiddleware("block", func(ctx *MsgContext, db NoData) bool {
		return false
	}))
	plugin.SetMessageFallback(func(ctx *MsgContext, db NoData) error {
		called = true
		return nil
	})

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
	}
	bot.AddPlugins(plugin)

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 7,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 1,
			Text:      "blocked",
			From:      &tgapi.User{ID: 41},
			Chat:      &tgapi.Chat{ID: 99},
		},
	})

	if called {
		t.Fatal("message fallback must not run when plugin middleware blocks")
	}
}

func TestHandleMessageFallbackDoesNotRunWhenCommandMatches(t *testing.T) {
	commandCalled := false
	fallbackCalled := false
	plugin := NewPlugin[NoData]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoData) error {
		commandCalled = true
		return nil
	}, "start")
	plugin.SetMessageFallback(func(ctx *MsgContext, db NoData) error {
		fallbackCalled = true
		return nil
	})

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
	}
	bot.AddPlugins(plugin)

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 8,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 1,
			Text:      "/start",
			From:      &tgapi.User{ID: 41},
			Chat:      &tgapi.Chat{ID: 99},
		},
	})

	if !commandCalled {
		t.Fatal("expected command handler to be called")
	}
	if fallbackCalled {
		t.Fatal("message fallback must not run when command matches")
	}
}

func TestHandleChannelPostCommandWithSenderChat(t *testing.T) {
	called := false
	plugin := NewPlugin[NoData]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoData) error {
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

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoData]{clonePlugin(plugin)},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 10,
		Type:     tgapi.UpdateTypeChannelPost,
		ChannelPost: &tgapi.Message{
			MessageID:  55,
			Text:       "/ping",
			SenderChat: &tgapi.Chat{ID: -1001, Type: tgapi.ChatTypeChannel},
			Chat:       &tgapi.Chat{ID: -1001, Type: tgapi.ChatTypeChannel},
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
	plugin := NewPlugin[NoData]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoData) error {
		return ctx.BindArgs(&got)
	}, "ban",
		NewCommandArg("user_id").SetValueType(CommandValueIntType).SetRequired(),
		NewCommandArg("reason").SetRequired(),
	)

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoData]{clonePlugin(plugin)},
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 11,
		Type:     tgapi.UpdateTypeMessage,
		Message: &tgapi.Message{
			MessageID: 1,
			Text:      "/ban 42 too loud",
			Chat:      &tgapi.Chat{ID: 99, Type: tgapi.ChatTypePrivate},
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
	plugin := NewPlugin[NoData]("test")
	plugin.NewPayload(func(ctx *MsgContext, db NoData) error {
		return ctx.BindArgs(&got)
	}, "approve",
		NewCommandArg("id").SetValueType(CommandValueIntType).SetRequired(),
		NewCommandArg("note").SetRequired(),
	)

	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		payloadType: BotPayloadJSON,
		plugins:     []Plugin[NoData]{clonePlugin(plugin)},
	}

	data, err := encodeJSONPayload(CallbackData{
		Command: "approve",
		Args:    []string{"7", "looks", "good"},
	})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
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

	plugin := NewPlugin[NoData]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoData) error {
		commandCalled = true
		return nil
	}, "ping")
	plugin.AddUpdateHandler(tgapi.UpdateTypeEditedMessage, func(ctx *MsgContext, db NoData) error {
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

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoData]{clonePlugin(plugin)},
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

	plugin := NewPlugin[NoData]("test")
	plugin.NewCommand(func(ctx *MsgContext, db NoData) error {
		commandCalled = true
		return nil
	}, "ping")
	plugin.AddUpdateHandler(tgapi.UpdateTypeEditedChannelPost, func(ctx *MsgContext, db NoData) error {
		updateCalled = true
		if ctx.Msg == nil {
			t.Fatal("expected ctx.Msg in edited channel post handler")
		}
		return nil
	})

	bot := &Bot[NoData]{
		logger:   sneklog.NewLogger(),
		prefixes: []string{"/"},
		plugins:  []Plugin[NoData]{clonePlugin(plugin)},
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
	plugin := NewPlugin[NoData]("test")
	plugin.NewPayload(func(ctx *MsgContext, db NoData) error {
		called = true
		if ctx.CallbackQueryID != "cb-msg" {
			t.Fatalf("unexpected CallbackQueryID: %q", ctx.CallbackQueryID)
		}
		if ctx.CallbackMsgID != 55 {
			t.Fatalf("unexpected CallbackMsgID: %d", ctx.CallbackMsgID)
		}
		if ctx.InlineMsgID != "" {
			t.Fatalf("did not expect InlineMsgID, got %q", ctx.InlineMsgID)
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

	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		payloadType: BotPayloadJSON,
		plugins:     []Plugin[NoData]{clonePlugin(plugin)},
	}

	data, err := encodeJSONPayload(CallbackData{Command: "approve", Args: []string{"7", "ok"}})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
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
	plugin := NewPlugin[NoData]("test")
	plugin.NewPayload(func(ctx *MsgContext, db NoData) error {
		called = true
		if ctx.CallbackQueryID != "cb-inline" {
			t.Fatalf("unexpected CallbackQueryID: %q", ctx.CallbackQueryID)
		}
		if ctx.CallbackMsgID != 0 {
			t.Fatalf("did not expect CallbackMsgID, got %d", ctx.CallbackMsgID)
		}
		if ctx.InlineMsgID != "inline-55" {
			t.Fatalf("unexpected InlineMsgID: %q", ctx.InlineMsgID)
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

	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		payloadType: BotPayloadJSON,
		plugins:     []Plugin[NoData]{clonePlugin(plugin)},
	}

	data, err := encodeJSONPayload(CallbackData{Command: "inline.approve", Args: []string{"9"}})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 31,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:              "cb-inline",
			Data:            data,
			From:            tgapi.User{ID: 8},
			InlineMessageID: new("inline-55"),
		},
	})

	if !called {
		t.Fatal("expected inline payload handler to be called")
	}
}

func TestHandleCallbackObserverEmitsPayloadEvents(t *testing.T) {
	observer := &recordingObserver{}
	plugin := NewPlugin[NoData]("test")
	plugin.NewPayload(func(ctx *MsgContext, db NoData) error {
		return nil
	}, "approve")

	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		payloadType: BotPayloadJSON,
		plugins:     []Plugin[NoData]{clonePlugin(plugin)},
		observer:    observer,
	}

	data, err := encodeJSONPayload(CallbackData{Command: "approve", Args: []string{"7"}})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 32,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-observer",
			Data: data,
			From: tgapi.User{ID: 7},
			Message: &tgapi.Message{
				MessageID: 56,
				Chat:      &tgapi.Chat{ID: 78},
			},
		},
	})

	if len(observer.started) != 1 {
		t.Fatalf("expected one started event, got %d", len(observer.started))
	}
	if got := observer.started[0]; got.HandlerKind != HandlerPayloadKind || got.HandlerName != "approve" || got.Plugin != "test" {
		t.Fatalf("unexpected started event: %#v", got)
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; got.HandlerKind != HandlerPayloadKind || got.HandlerName != "approve" || got.Plugin != "test" || got.Err != nil || got.UserFacing {
		t.Fatalf("unexpected finished event: %#v", got)
	}
	if len(observer.errors) != 0 {
		t.Fatalf("did not expect error events, got %#v", observer.errors)
	}
}

func TestHandleCallbackObserverEmitsPayloadErrors(t *testing.T) {
	observer := &recordingObserver{}
	plugin := NewPlugin[NoData]("test")
	wantErr := AsInternalError(errors.New("boom"))
	plugin.NewPayload(func(ctx *MsgContext, db NoData) error {
		return wantErr
	}, "approve")

	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		payloadType: BotPayloadJSON,
		plugins:     []Plugin[NoData]{clonePlugin(plugin)},
		observer:    observer,
	}

	data, err := encodeJSONPayload(CallbackData{Command: "approve", Args: []string{"7"}})
	if err != nil {
		t.Fatalf("encodeJSONPayload returned error: %v", err)
	}

	bot.handle(context.Background(), &tgapi.Update{
		UpdateID: 33,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-observer-err",
			Data: data,
			From: tgapi.User{ID: 7},
			Message: &tgapi.Message{
				MessageID: 57,
				Chat:      &tgapi.Chat{ID: 79},
			},
		},
	})

	if len(observer.started) != 1 {
		t.Fatalf("expected one started event, got %d", len(observer.started))
	}
	if len(observer.finished) != 1 {
		t.Fatalf("expected one finished event, got %d", len(observer.finished))
	}
	if got := observer.finished[0]; !errors.Is(got.Err, wantErr) || got.UserFacing {
		t.Fatalf("unexpected finished event: %#v", got)
	}
	if len(observer.errors) != 1 {
		t.Fatalf("expected one error event, got %d", len(observer.errors))
	}
	if got := observer.errors[0]; !errors.Is(got.Err, wantErr) || got.HandlerKind != HandlerPayloadKind || got.HandlerName != "approve" || got.Plugin != "test" || got.UserFacing {
		t.Fatalf("unexpected error event: %#v", got)
	}
}

func TestHandleCallbackObserverEmitsDecodeErrors(t *testing.T) {
	observer := &recordingObserver{}
	bot := &Bot[NoData]{
		logger:      sneklog.NewLogger(),
		payloadType: BotPayloadJSON,
		observer:    observer,
	}

	handled := bot.handleCallback(&tgapi.Update{
		UpdateID: 34,
		Type:     tgapi.UpdateTypeCallbackQuery,
		CallbackQuery: &tgapi.CallbackQuery{
			ID:   "cb-bad",
			Data: "{not-json",
			From: tgapi.User{ID: 7},
		},
	}, &MsgContext{
		Update: tgapi.Update{
			UpdateID: 34,
			Type:     tgapi.UpdateTypeCallbackQuery,
		},
		Logger:          bot.logger,
		ctx:             context.Background(),
		CallbackQueryID: "cb-bad",
		From:            &tgapi.User{ID: 7},
		FromID:          7,
		sceneRuntime:    bot,
	})

	if handled {
		t.Fatal("expected invalid callback payload to stay unhandled")
	}
	if len(observer.started) != 0 || len(observer.finished) != 0 {
		t.Fatalf("expected no handler lifecycle events for decode failure, got started=%d finished=%d", len(observer.started), len(observer.finished))
	}
	if len(observer.errors) != 1 {
		t.Fatalf("expected one observer error event, got %d", len(observer.errors))
	}
	ev := observer.errors[0]
	if ev.Plugin != "bot" {
		t.Fatalf("unexpected plugin: %q", ev.Plugin)
	}
	if ev.HandlerKind != HandlerPayloadKind {
		t.Fatalf("unexpected handler kind: %q", ev.HandlerKind)
	}
	if ev.HandlerName != "decodePayload" {
		t.Fatalf("unexpected handler name: %q", ev.HandlerName)
	}
	if ev.UserFacing {
		t.Fatal("expected decode failure to stay internal")
	}
}
