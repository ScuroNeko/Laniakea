package laniakea

import (
	"errors"
	"strings"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
)

func TestDraftFlushRequiresChatID(t *testing.T) {
	draft := NewRandomDraftProvider(&tgapi.API{}).NewDraft(tgapi.ParseNone)
	draft.Message = "hello"

	if err := draft.Flush(); err != ErrDraftChatIDZero {
		t.Fatalf("expected ErrDraftChatIDZero, got %v", err)
	}
}

func TestMsgContextNewDraftWorksWithoutLimiter(t *testing.T) {
	ctx := &MsgContext{
		Api: &tgapi.API{},
		Msg: &tgapi.Message{
			Chat: &tgapi.Chat{ID: 42, Type: string(tgapi.ChatTypePrivate)},
		},
		Logger:        slog.CreateLogger(),
		draftProvider: NewRandomDraftProvider(&tgapi.API{}),
	}

	draft := ctx.NewDraft()
	if draft == nil {
		t.Fatal("expected draft")
	}
	if draft.chatID != 42 {
		t.Fatalf("unexpected chat id: %d", draft.chatID)
	}
}

func TestDraftFlushRejectsLongMessage(t *testing.T) {
	draft := NewRandomDraftProvider(&tgapi.API{}).NewDraft(tgapi.ParseNone).SetChat(42, 0)
	draft.Message = strings.Repeat("a", maxMessageTextLen+1)

	if err := draft.Flush(); !errors.Is(err, ErrMessageTooLong) {
		t.Fatalf("expected ErrMessageTooLong, got %v", err)
	}
}

func TestDraftPushRejectsLongMessage(t *testing.T) {
	draft := NewRandomDraftProvider(&tgapi.API{}).NewDraft(tgapi.ParseNone).SetChat(42, 0)

	if err := draft.Push(strings.Repeat("a", maxMessageTextLen+1)); !errors.Is(err, ErrMessageTooLong) {
		t.Fatalf("expected ErrMessageTooLong, got %v", err)
	}
}
