package laniakea

import (
	"strings"
	"testing"
)

func TestInlineKeyboardWrapsRowsAndEncodesJSONPayloads(t *testing.T) {
	kb := NewInlineKeyboardJson(2).
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
				SetUrl("https://example.test"),
		)

	button := kb.Get().InlineKeyboard[0][0]
	if button.Style != ButtonStylePrimary {
		t.Fatalf("unexpected style: %q", button.Style)
	}
	if button.URL != "https://example.test" {
		t.Fatalf("unexpected url: %q", button.URL)
	}
}
