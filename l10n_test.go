package laniakea

import (
	"fmt"
	"sync"
	"testing"
)

func TestL10nTranslateUsesFallbackAndKey(t *testing.T) {
	l10n := NewL10n("en").
		AddDictEntry("greeting", DictEntry{"en": "Hello", "ru": "Privet"}).
		AddDictEntry("partial", DictEntry{"ru": "Tolko ru"})

	tests := []struct {
		name string
		lang string
		key  string
		want string
	}{
		{name: "exact match", lang: "ru", key: "greeting", want: "Privet"},
		{name: "fallback language", lang: "es", key: "greeting", want: "Hello"},
		{name: "missing fallback returns key", lang: "en", key: "partial", want: "partial"},
		{name: "unknown key returns key", lang: "en", key: "unknown", want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := l10n.Translate(tt.lang, tt.key); got != tt.want {
				t.Fatalf("unexpected translation: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestL10nAddDictEntryCopiesInput(t *testing.T) {
	l10n := NewL10n("en")
	entry := DictEntry{"en": "Hello"}

	l10n.AddDictEntry("greeting", entry)
	entry["en"] = "Mutated"

	if got := l10n.Translate("en", "greeting"); got != "Hello" {
		t.Fatalf("unexpected translation after external mutation: got %q", got)
	}
}

func TestL10nZeroValueIsUsable(t *testing.T) {
	var l10n L10n

	l10n.AddDictEntry("greeting", DictEntry{"en": "Hello"})

	if got := l10n.Translate("en", "greeting"); got != "Hello" {
		t.Fatalf("unexpected translation from zero-value l10n: got %q", got)
	}
}

func TestL10nConcurrentAccess(t *testing.T) {
	l10n := NewL10n("en")
	l10n.AddDictEntry("base", DictEntry{"en": "Hello"})

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				l10n.AddDictEntry(fmt.Sprintf("key-%d-%d", i, j), DictEntry{"en": "value"})
				_ = l10n.Translate("en", "base")
			}
		}(i)
	}
	wg.Wait()

	if got := l10n.Translate("en", "base"); got != "Hello" {
		t.Fatalf("unexpected translation after concurrent access: got %q", got)
	}
}
