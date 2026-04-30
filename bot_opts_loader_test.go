package laniakea

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

func TestBotOptsFileJSONCodecRoundTrip(t *testing.T) {
	codec := BotOptsFileJSONCodec{}
	want := &BotOpts{
		Token:                 "TOKEN",
		UpdateTypes:           []tgapi.UpdateType{tgapi.UpdateTypeMessage, tgapi.UpdateTypeCallbackQuery},
		Debug:                 true,
		ErrorTemplate:         "Error: %s",
		Prefixes:              []string{"/", "!"},
		LoggerBasePath:        "/tmp/logs",
		UseRequestLogger:      true,
		WriteToFile:           true,
		UseTestServer:         true,
		APIURL:                "https://api.example.invalid",
		RateLimit:             42,
		DropRateLimitOverflow: true,
		StrictPayloadType:     true,
		MaxWorkers:            64,
		FileConfigVersion:     ConfigVersion,
	}

	data, err := codec.ToBytes(want)
	if err != nil {
		t.Fatalf("ToBytes returned error: %v", err)
	}

	got, err := codec.FromBytes(data)
	if err != nil {
		t.Fatalf("FromBytes returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestLoadBotOptsFileExpandsEnvPlaceholders(t *testing.T) {
	t.Setenv("TG_TOKEN", "TOKEN_FROM_ENV")
	t.Setenv("BOT_API_URL", "https://api.example.invalid")

	dir := t.TempDir()
	filename := filepath.Join(dir, "config.json")
	data := []byte(`{
		"token": "{{ TG_TOKEN }}",
		"api": {
			"url": "{{BOT_API_URL}}"
		},
		"error_template": "Error: %s"
	}`)
	if err := os.WriteFile(filename, data, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	got, err := LoadBotOptsFile(BotOptsFileJSONCodec{}, filename)
	if err != nil {
		t.Fatalf("LoadBotOptsFile returned error: %v", err)
	}

	if got.Token != "TOKEN_FROM_ENV" {
		t.Fatalf("unexpected token: got %q want %q", got.Token, "TOKEN_FROM_ENV")
	}
	if got.APIURL != "https://api.example.invalid" {
		t.Fatalf("unexpected api url: got %q want %q", got.APIURL, "https://api.example.invalid")
	}
	if got.ErrorTemplate != "Error: %s" {
		t.Fatalf("unexpected error template: got %q", got.ErrorTemplate)
	}
	if got.FileConfigVersion != 0 {
		t.Fatalf("unexpected file config version: got %d want 0", got.FileConfigVersion)
	}
}

func TestLoadBotOptsFileReturnsDecodeError(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "config.json")
	if err := os.WriteFile(filename, []byte(`{"token":`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if _, err := LoadBotOptsFile(BotOptsFileJSONCodec{}, filename); err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestSaveBotOptsFileWritesEncodedData(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "config.json")
	want := &BotOpts{
		Token:             "TOKEN",
		UpdateTypes:       []tgapi.UpdateType{tgapi.UpdateTypeMessage},
		ErrorTemplate:     "Error: %s",
		Prefixes:          []string{"/"},
		APIURL:            "https://api.example.invalid",
		RateLimit:         30,
		MaxWorkers:        32,
		FileConfigVersion: ConfigVersion,
	}

	if err := SaveBotOptsFile(BotOptsFileJSONCodec{}, filename, want); err != nil {
		t.Fatalf("SaveBotOptsFile returned error: %v", err)
	}

	got, err := LoadBotOptsFile(BotOptsFileJSONCodec{}, filename)
	if err != nil {
		t.Fatalf("LoadBotOptsFile returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("saved file mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestLoadBotOptsFileRejectsFutureConfigVersion(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "config.json")
	data := []byte(`{
		"version": 2,
		"token": "TOKEN"
	}`)
	if err := os.WriteFile(filename, data, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	_, err := LoadBotOptsFile(BotOptsFileJSONCodec{}, filename)
	if !errors.Is(err, ErrConfigVersionMismatch) {
		t.Fatalf("expected ErrConfigVersionMismatch, got %v", err)
	}
}
