package laniakea

import (
	"testing"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/slog"
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
