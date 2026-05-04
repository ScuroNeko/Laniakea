package tgfmt

import "testing"

func TestEscapeMarkdownV2(t *testing.T) {
	got := EscapeMarkdownV2(`a_b*c[1](x)!`)
	want := MarkdownV2(`a\_b\*c\[1\]\(x\)\!`)

	if got != want {
		t.Fatalf("EscapeMarkdownV2() = %q, want %q", got, want)
	}
}

func TestMarkdownV2ComposesWithoutDoubleEscaping(t *testing.T) {
	got := EscapeMarkdownV2("a*b").Bold().Italic()
	want := MarkdownV2(`_*a\*b*_`)

	if got != want {
		t.Fatalf("formatted text = %q, want %q", got, want)
	}
}

func TestMarkdownV2FormattingMethods(t *testing.T) {
	tests := []struct {
		name string
		got  MarkdownV2
		want MarkdownV2
	}{
		{name: "bold", got: EscapeMarkdownV2("text").Bold(), want: "*text*"},
		{name: "italic", got: EscapeMarkdownV2("text").Italic(), want: "_text_"},
		{name: "underline", got: EscapeMarkdownV2("text").Underline(), want: "__text__"},
		{name: "strikethrough", got: EscapeMarkdownV2("text").Strikethrough(), want: "~text~"},
		{name: "spoiler", got: EscapeMarkdownV2("text").Spoiler(), want: "||text||"},
		{name: "inline code", got: EscapeMarkdownV2("text").InlineCode(), want: "`text`"},
		{name: "block code", got: EscapeMarkdownV2("text").BlockCode(), want: "```\ntext\n```"},
		{name: "block code language", got: EscapeMarkdownV2("text").BlockCodeLanguage("go"), want: "```go\ntext\n```"},
		{name: "quote", got: EscapeMarkdownV2("a\nb").Quote(), want: ">a\n>b"},
		{name: "expandable quote", got: EscapeMarkdownV2("text").QuoteExpandable(), want: "**>text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("formatted text = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestMarkdownV2LinkEscapesDestination(t *testing.T) {
	got := EscapeMarkdownV2("Laniakea").Link(`https://example.test/a)b\c`)
	want := MarkdownV2(`[Laniakea](https://example.test/a\)b\\c)`)

	if got != want {
		t.Fatalf("Link() = %q, want %q", got, want)
	}
}

func TestMarkdownV2SpecialLinks(t *testing.T) {
	tests := []struct {
		name string
		got  MarkdownV2
		want MarkdownV2
	}{
		{name: "mention", got: EscapeMarkdownV2("User").Mention(123), want: "[User](tg://user?id=123)"},
		{name: "emoji", got: EscapeMarkdownV2("emoji").Emoji(`12)3`), want: `[emoji](tg://emoji?id=12\)3)`},
		{name: "time", got: EscapeMarkdownV2("date").Time(1772323200), want: "![date](tg://time?unix=1772323200)"},
		{name: "time format", got: EscapeMarkdownV2("date").TimeFormat(1772323200, `MMM ) yyyy`), want: `![date](tg://time?unix=1772323200&format=MMM \) yyyy)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("formatted link = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
