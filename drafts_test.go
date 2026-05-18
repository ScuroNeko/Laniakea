package laniakea

import (
	"errors"
	"strings"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

func TestDraftFlushRequiresChatID(t *testing.T) {
	draft := NewRandomDraftProvider(&tgapi.API{}).NewDraft(tgapi.ParseNone)
	draft.Message = "hello"

	if err := draft.Flush(); !errors.Is(err, ErrDraftChatIDZero) {
		t.Fatalf("expected ErrDraftChatIDZero, got %v", err)
	}
}

func TestMsgContextNewDraftWorksWithoutLimiter(t *testing.T) {
	ctx := &MessageContext{
		API: &tgapi.API{},
		Msg: &tgapi.Message{
			Chat: &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate},
		},
		Logger:        sneklog.NewLogger(),
		draftProvider: NewRandomDraftProvider(&tgapi.API{}),
	}

	draft := ctx.NewDraft()
	if draft == nil {
		t.Fatal("expected draft")
		return
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

// TestDraftPushLeavesMessageUnchangedOnValidationFailure covers the validation
// order fix: when the candidate Message (current + new text) overflows the
// Telegram limit, the existing Message must remain intact so callers can
// recover and retry with a shorter payload instead of finding the draft in
// a half-mutated state.
func TestDraftPushLeavesMessageUnchangedOnValidationFailure(t *testing.T) {
	draft := NewRandomDraftProvider(&tgapi.API{}).NewDraft(tgapi.ParseNone).SetChat(42, 0)
	draft.Message = "hello"

	if err := draft.Push(strings.Repeat("a", maxMessageTextLen+1)); !errors.Is(err, ErrMessageTooLong) {
		t.Fatalf("expected ErrMessageTooLong, got %v", err)
	}
	if draft.Message != "hello" {
		t.Fatalf("expected draft Message to stay %q, got %q", "hello", draft.Message)
	}
}
