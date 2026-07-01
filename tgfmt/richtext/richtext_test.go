package richtext

import (
	"encoding/json"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

func roundtrip(t *testing.T, in RichText) {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out, err := Unmarshal(b)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", b, err)
	}
	b2, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("remarshal: %v", err)
	}
	if string(b) != string(b2) {
		t.Fatalf("not stable:\n  %s\n  %s", b, b2)
	}
}

func TestRoundtrip(t *testing.T) {
	cases := []RichText{
		String("hello"),
		Array{String("a "), Bold(String("b")), String(" c")},
		Bold(Italic(String("nested"))),
		URL{String("Anthropic"), "https://anthropic.com"},
		CustomEmoji{"5368324170671202286", "👍"},
		MathematicalExpression{"x^2 + y^2"},
		Anchor{"chapter-1"},
		DateTime{String("22:45 tomorrow"), 1647531900, "wDT"},
		TextMention{String("Bob"), tgapi.User{ID: 42, FirstName: "Bob"}},
		AnchorLink{String("back to top"), ""},
		Reference{String("ref"), "note-1"},
		// глубокая вложенность
		Bold(Array{
			String("bold and "),
			Italic(Underline(String("deep"))),
			Spoiler(CustomEmoji{"1", "x"}),
		}),
	}
	for _, c := range cases {
		roundtrip(t, c)
	}
}

func TestPlainFormsAreBare(t *testing.T) {
	b, _ := json.Marshal(String("hi"))
	if string(b) != `"hi"` {
		t.Fatalf("string should be bare: %s", b)
	}
	b, _ = json.Marshal(Array{String("a"), String("b")})
	if string(b) != `["a","b"]` {
		t.Fatalf("array should be bare: %s", b)
	}
}

func TestLeafHasNoText(t *testing.T) {
	b, _ := json.Marshal(Anchor{"x"})
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if _, ok := m["text"]; ok {
		t.Fatalf("anchor must not have text field: %s", b)
	}
}
