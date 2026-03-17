package laniakea

import "testing"

func TestCheckPrefixesSkipsEmptyPrefixes(t *testing.T) {
	bot := &Bot[NoDB]{prefixes: []string{"", "/"}}

	if prefix, ok := bot.checkPrefixes("hello"); ok {
		t.Fatalf("unexpected prefix match for plain text: %q", prefix)
	}
	if prefix, ok := bot.checkPrefixes("/start"); !ok || prefix != "/" {
		t.Fatalf("unexpected prefix result: prefix=%q ok=%v", prefix, ok)
	}
}
