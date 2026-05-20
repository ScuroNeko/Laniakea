package tgfmt

import "testing"

func TestEscapeMarkdown(t *testing.T) {
	got := EscapeMarkdown("a_b*c[1]`x`")
	want := Markdown("a\\_b\\*c\\[1]\\`x\\`")

	if got != want {
		t.Fatalf("EscapeMarkdown() = %q, want %q", got, want)
	}
}

func TestMarkdownFormattingMethods(t *testing.T) {
	tests := []struct {
		name string
		got  Markdown
		want Markdown
	}{
		{name: "bold", got: EscapeMarkdown("text").Bold(), want: "*text*"},
		{name: "italic", got: EscapeMarkdown("text").Italic(), want: "_text_"},
		{name: "link", got: EscapeMarkdown("Laniakea").Link("https://example.test"), want: "[Laniakea](https://example.test)"},
		{name: "mention", got: EscapeMarkdown("User").Mention(123), want: "[User](tg://user?id=123)"},
		{name: "inline code", got: EscapeMarkdown("text").InlineCode(), want: "`text`"},
		{name: "block code", got: EscapeMarkdown("text").BlockCode(), want: "```\ntext\n```"},
		{name: "block code language", got: EscapeMarkdown("text").BlockCodeLanguage("go"), want: "```go\ntext\n```"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("formatted Markdown = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
