package laniakea

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestInlineKeyboardWrapsRowsAndEncodesJSONPayloads(t *testing.T) {
	kb := NewInlineKeyboardJSON(2).
		AddCallbackButton("A", "cmd", 1).
		AddCallbackButton("B", "cmd", 2).
		AddCallbackButton("C", "cmd", 3)

	markup := kb.Get()
	if got := len(markup.InlineKeyboard); got != 2 {
		t.Fatalf("unexpected row count: %d", got)
	}
	if got := len(markup.InlineKeyboard[0]); got != 2 {
		t.Fatalf("unexpected first row size: %d", got)
	}
	if got := len(markup.InlineKeyboard[1]); got != 1 {
		t.Fatalf("unexpected second row size: %d", got)
	}
	if !strings.Contains(markup.InlineKeyboard[0][0].CallbackData, `"cmd":"cmd"`) {
		t.Fatalf("expected JSON callback payload, got %q", markup.InlineKeyboard[0][0].CallbackData)
	}
}

func TestInlineKeyboardBuilderPreservesConfiguredButtonFields(t *testing.T) {
	kb := NewInlineKeyboardBase64(3).
		AddButton(
			NewInlineKbButton("Docs").
				SetStyle(ButtonStylePrimary).
				SetURL("https://example.test"),
		)

	button := kb.Get().InlineKeyboard[0][0]
	if button.Style != ButtonStylePrimary {
		t.Fatalf("unexpected style: %q", button.Style)
	}
	if button.URL != "https://example.test" {
		t.Fatalf("unexpected url: %q", button.URL)
	}
}

func TestInlineKeyboardGetPayloadTypeReturnsLocalOverride(t *testing.T) {
	kb := NewInlineKeyboardJSON(2)
	if got := kb.GetPayloadType(); got != BotPayloadJSON {
		t.Fatalf("unexpected initial payload type: %q", got)
	}
	kb.SetPayloadType(BotPayloadBase64)
	if got := kb.GetPayloadType(); got != BotPayloadBase64 {
		t.Fatalf("unexpected updated payload type: %q", got)
	}
}

func TestDecodePayloadAcceptsBase64KeyboardPayloadWhenBotPrefersJSON(t *testing.T) {
	kb := NewInlineKeyboardBase64(1).
		AddCallbackButton("A", "cmd", 1, "two")

	got, _, err := decodePayload(BotPayloadJSON, kb.Get().InlineKeyboard[0][0].CallbackData, false)
	if err != nil {
		t.Fatalf("decodePayload returned error: %v", err)
	}

	want := CallbackData{Command: "cmd", Args: []string{"1", "two"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected payload: got %#v want %#v", got, want)
	}
}

func TestDecodePayloadAcceptsJSONKeyboardPayloadWhenBotPrefersBase64(t *testing.T) {
	kb := NewInlineKeyboardJSON(1).
		AddCallbackButton("A", "cmd", 1, "two")

	got, _, err := decodePayload(BotPayloadBase64, kb.Get().InlineKeyboard[0][0].CallbackData, false)
	if err != nil {
		t.Fatalf("decodePayload returned error: %v", err)
	}

	want := CallbackData{Command: "cmd", Args: []string{"1", "two"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected payload: got %#v want %#v", got, want)
	}
}

func TestDecodePayloadStrictRejectsMismatchedType(t *testing.T) {
	kb := NewInlineKeyboardBase64(1).
		AddCallbackButton("A", "cmd", 1)

	_, _, err := decodePayload(BotPayloadJSON, kb.Get().InlineKeyboard[0][0].CallbackData, true)
	if !errors.Is(err, ErrPayloadTypeMismatch) {
		t.Fatalf("expected ErrPayloadTypeMismatch, got %v", err)
	}
}
