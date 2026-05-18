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
			NewInlineKeyboardButton("Docs").
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

func TestInlineKeyboardButtonBuilderSetCallbackDataDefaultsToJSON(t *testing.T) {
	kb := NewInlineKeyboardBase64(1).
		AddButton(NewInlineKeyboardButton("A").SetCallbackData("cmd", 1, "two"))

	button := kb.Get().InlineKeyboard[0][0]
	if !strings.Contains(button.CallbackData, `"cmd":"cmd"`) {
		t.Fatalf("expected JSON callback payload, got %q", button.CallbackData)
	}
}

func TestInlineKeyboardButtonBuilderSetCallbackDataUsesConfiguredPayloadType(t *testing.T) {
	kb := NewInlineKeyboardJSON(1).
		AddButton(NewInlineKeyboardButton("A").
			SetPayloadType(BotPayloadBase64).
			SetCallbackData("cmd", 1, "two"),
		)

	got, _, err := decodePayload(BotPayloadJSON, kb.Get().InlineKeyboard[0][0].CallbackData, false)
	if err != nil {
		t.Fatalf("decodePayload returned error: %v", err)
	}

	want := CallbackData{Command: "cmd", Args: []string{"1", "two"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected payload: got %#v want %#v", got, want)
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

func TestDecodePayloadAcceptsCompactKeyboardPayloadWhenBotPrefersJSON(t *testing.T) {
	kb := NewInlineKeyboardCompact(1).
		AddCallbackButton("A", "cmd", 1, "two")

	got, decodedType, err := decodePayload(BotPayloadJSON, kb.Get().InlineKeyboard[0][0].CallbackData, false)
	if err != nil {
		t.Fatalf("decodePayload returned error: %v", err)
	}
	if decodedType != BotPayloadCompact {
		t.Fatalf("unexpected decoded payload type: got %q want %q", decodedType, BotPayloadCompact)
	}

	want := CallbackData{Command: "cmd", Args: []string{"1", "two"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected payload: got %#v want %#v", got, want)
	}
}

func TestDecodePayloadAcceptsCompactBase64KeyboardPayloadWhenBotPrefersJSON(t *testing.T) {
	kb := NewInlineKeyboardCompactBase64(1).
		AddCallbackButton("A", "cmd", 1, "two")

	got, decodedType, err := decodePayload(BotPayloadJSON, kb.Get().InlineKeyboard[0][0].CallbackData, false)
	if err != nil {
		t.Fatalf("decodePayload returned error: %v", err)
	}
	if decodedType != BotPayloadCompactBase64 {
		t.Fatalf("unexpected decoded payload type: got %q want %q", decodedType, BotPayloadCompactBase64)
	}

	want := CallbackData{Command: "cmd", Args: []string{"1", "two"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected payload: got %#v want %#v", got, want)
	}
}

// TestCompactPayloadRoundTripsWithSeparatorChars guards the compact-encoding
// escape fix. Args containing the , | or \ separator bytes previously corrupted
// on decode; now they must round-trip exactly.
//
// Note: the compact format coalesces "no args" with "single empty arg" — both
// emit "cmd|" and decode to nil args. Use other encodings if that distinction
// matters.
func TestCompactPayloadRoundTripsWithSeparatorChars(t *testing.T) {
	tests := []struct {
		name string
		data CallbackData
	}{
		{name: "plain", data: CallbackData{Command: "cmd", Args: []string{"one", "two"}}},
		{name: "no args", data: CallbackData{Command: "cmd"}},
		{name: "comma in arg", data: CallbackData{Command: "cmd", Args: []string{"a,b", "c"}}},
		{name: "pipe in arg", data: CallbackData{Command: "cmd", Args: []string{"a|b", "c"}}},
		{name: "backslash in arg", data: CallbackData{Command: "cmd", Args: []string{`a\b`, "c"}}},
		{name: "all specials in arg", data: CallbackData{Command: "cmd", Args: []string{`a,b|c\d`}}},
		{name: "specials in command", data: CallbackData{Command: "a|b,c", Args: []string{"x"}}},
		{name: "two empty args", data: CallbackData{Command: "cmd", Args: []string{"", ""}}},
		{name: "utf8 args", data: CallbackData{Command: "cmd", Args: []string{"привет", "мир"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := encodeCompactPayload(tt.data)
			if err != nil {
				t.Fatalf("encodeCompactPayload returned error: %v", err)
			}
			got, err := decodeCompactPayload(encoded)
			if err != nil {
				t.Fatalf("decodeCompactPayload returned error: %v", err)
			}
			if got.Command != tt.data.Command {
				t.Fatalf("command mismatch: got %q want %q (encoded=%q)", got.Command, tt.data.Command, encoded)
			}
			if len(got.Args) != len(tt.data.Args) {
				t.Fatalf("args length mismatch: got %v want %v (encoded=%q)", got.Args, tt.data.Args, encoded)
			}
			for i := range tt.data.Args {
				if got.Args[i] != tt.data.Args[i] {
					t.Fatalf("arg %d mismatch: got %q want %q (encoded=%q)", i, got.Args[i], tt.data.Args[i], encoded)
				}
			}
		})
	}
}

func TestCompactPayloadDecodeRejectsMissingSeparator(t *testing.T) {
	if _, err := decodeCompactPayload("noseparator"); err == nil {
		t.Fatal("expected error decoding payload without separator")
	}
}

func TestDecodePayloadStrictRejectsCompactMismatchedType(t *testing.T) {
	kb := NewInlineKeyboardCompact(1).
		AddCallbackButton("A", "cmd", 1)

	_, _, err := decodePayload(BotPayloadJSON, kb.Get().InlineKeyboard[0][0].CallbackData, true)
	if !errors.Is(err, ErrPayloadTypeMismatch) {
		t.Fatalf("expected ErrPayloadTypeMismatch, got %v", err)
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
