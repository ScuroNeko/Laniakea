package tgapi

import (
	"encoding/json"
	"testing"
)

func roundtripRichText(t *testing.T, in RichText) {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out, err := UnmarshalRichText(b)
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

func TestRichTextRoundtrip(t *testing.T) {
	cases := []RichText{
		RichTextPlain("hello"),
		RichTextArray{RichTextPlain("a "), RichTextWrap{"bold", RichTextPlain("b")}, RichTextPlain(" c")},
		RichTextWrap{"bold", RichTextWrap{"italic", RichTextPlain("nested")}},
		RichTextURL{RichTextPlain("Anthropic"), "https://anthropic.com"},
		RichTextCustomEmoji{"5368324170671202286", "👍"},
		RichTextMathematicalExpression{"x^2 + y^2"},
		RichTextAnchor{"chapter-1"},
		RichTextDateTime{RichTextPlain("22:45 tomorrow"), 1647531900, "wDT"},
		RichTextTextMention{RichTextPlain("Bob"), User{ID: 42, FirstName: "Bob"}},
		RichTextAnchorLink{RichTextPlain("back to top"), ""},
		RichTextReference{RichTextPlain("ref"), "note-1"},
		// deep nesting
		RichTextWrap{"bold", RichTextArray{
			RichTextPlain("bold and "),
			RichTextWrap{"italic", RichTextWrap{"underline", RichTextPlain("deep")}},
			RichTextWrap{"spoiler", RichTextCustomEmoji{"1", "x"}},
		}},
	}
	for _, c := range cases {
		roundtripRichText(t, c)
	}
}

func TestRichTextPlainFormsAreBare(t *testing.T) {
	b, _ := json.Marshal(RichTextPlain("hi"))
	if string(b) != `"hi"` {
		t.Fatalf("string should be bare: %s", b)
	}
	b, _ = json.Marshal(RichTextArray{RichTextPlain("a"), RichTextPlain("b")})
	if string(b) != `["a","b"]` {
		t.Fatalf("array should be bare: %s", b)
	}
}

func TestRichTextLeafHasNoText(t *testing.T) {
	b, _ := json.Marshal(RichTextAnchor{"x"})
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if _, ok := m["text"]; ok {
		t.Fatalf("anchor must not have text field: %s", b)
	}
}
