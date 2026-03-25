package tgapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseNoneOmitsParseModeInJSON(t *testing.T) {
	data, err := json.Marshal(SendMessageP{
		ChatID:    42,
		Text:      "hello",
		ParseMode: ParseNone,
	})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	if strings.Contains(string(data), `"parse_mode"`) {
		t.Fatalf("expected parse_mode to be omitted, got %s", string(data))
	}
}

func TestParseModeStillSerializesExplicitModes(t *testing.T) {
	data, err := json.Marshal(SendMessageP{
		ChatID:    42,
		Text:      "hello",
		ParseMode: ParseMDV2,
	})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	if !strings.Contains(string(data), `"parse_mode":"MarkdownV2"`) {
		t.Fatalf("expected MarkdownV2 parse_mode, got %s", string(data))
	}
}
