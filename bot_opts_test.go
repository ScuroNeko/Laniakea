package laniakea

import (
	"reflect"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

func TestLoadOptsFromEnvIgnoresEmptyUpdateTypes(t *testing.T) {
	t.Setenv("UPDATE_TYPES", "")

	opts := LoadOptsFromEnv()
	if len(opts.UpdateTypes) != 0 {
		t.Fatalf("expected no update types, got %v", opts.UpdateTypes)
	}
}

func TestLoadOptsFromEnvSplitsAndTrimsUpdateTypes(t *testing.T) {
	t.Setenv("UPDATE_TYPES", "message; ; callback_query ")

	opts := LoadOptsFromEnv()
	want := []tgapi.UpdateType{tgapi.UpdateTypeMessage, tgapi.UpdateTypeCallbackQuery}
	if !reflect.DeepEqual(opts.UpdateTypes, want) {
		t.Fatalf("unexpected update types: got %v want %v", opts.UpdateTypes, want)
	}
}

func TestLoadPrefixesFromEnvDefaultsOnEmptyValue(t *testing.T) {
	t.Setenv("PREFIXES", "")

	got := LoadPrefixesFromEnv()
	want := []string{"/"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected prefixes: got %v want %v", got, want)
	}
}

func TestLoadPrefixesFromEnvDropsEmptyValues(t *testing.T) {
	t.Setenv("PREFIXES", "/; ; ! ")

	got := LoadPrefixesFromEnv()
	want := []string{"/", "!"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected prefixes: got %v want %v", got, want)
	}
}

func TestLoadOptsFromEnvReadsStrictPayloadType(t *testing.T) {
	t.Setenv("STRICT_PAYLOAD_TYPE", "true")

	opts := LoadOptsFromEnv()
	if !opts.StrictPayloadType {
		t.Fatal("expected StrictPayloadType to be enabled")
	}
}
