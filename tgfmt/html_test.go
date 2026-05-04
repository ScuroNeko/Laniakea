package tgfmt

import "testing"

func TestEscapeHTML(t *testing.T) {
	got := EscapeHTML(`<tag attr="a&b">`)
	want := HTML(`&lt;tag attr=&quot;a&amp;b&quot;&gt;`)

	if got != want {
		t.Fatalf("EscapeHTML() = %q, want %q", got, want)
	}
}

func TestHTMLComposesWithoutDoubleEscaping(t *testing.T) {
	got := EscapeHTML("<b>").Bold().Italic()
	want := HTML("<i><b>&lt;b&gt;</b></i>")

	if got != want {
		t.Fatalf("formatted HTML = %q, want %q", got, want)
	}
}

func TestHTMLFormattingMethods(t *testing.T) {
	tests := []struct {
		name string
		got  HTML
		want HTML
	}{
		{name: "bold", got: EscapeHTML("text").Bold(), want: "<b>text</b>"},
		{name: "italic", got: EscapeHTML("text").Italic(), want: "<i>text</i>"},
		{name: "underline", got: EscapeHTML("text").Underline(), want: "<u>text</u>"},
		{name: "strikethrough", got: EscapeHTML("text").Strikethrough(), want: "<s>text</s>"},
		{name: "spoiler", got: EscapeHTML("text").Spoiler(), want: "<tg-spoiler>text</tg-spoiler>"},
		{name: "inline code", got: EscapeHTML("text").InlineCode(), want: "<code>text</code>"},
		{name: "block code", got: EscapeHTML("text").BlockCode(), want: "<pre>text</pre>"},
		{name: "block code language", got: EscapeHTML("text").BlockCodeLanguage(`go"`), want: `<pre><code class="language-go&quot;">text</code></pre>`},
		{name: "quote", got: EscapeHTML("text").Quote(), want: "<blockquote>text</blockquote>"},
		{name: "expandable quote", got: EscapeHTML("text").QuoteExpandable(), want: "<blockquote expandable>text</blockquote>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("formatted HTML = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestHTMLLinkEscapesAttributes(t *testing.T) {
	got := EscapeHTML("Laniakea").Link(`https://example.test/?q="a&b"`)
	want := HTML(`<a href="https://example.test/?q=&quot;a&amp;b&quot;">Laniakea</a>`)

	if got != want {
		t.Fatalf("Link() = %q, want %q", got, want)
	}
}

func TestHTMLSpecialLinks(t *testing.T) {
	tests := []struct {
		name string
		got  HTML
		want HTML
	}{
		{name: "mention", got: EscapeHTML("User").Mention(123), want: `<a href="tg://user?id=123">User</a>`},
		{name: "emoji", got: EscapeHTML("emoji").Emoji(`12"3`), want: `<tg-emoji emoji-id="12&quot;3">emoji</tg-emoji>`},
		{name: "time", got: EscapeHTML("date").Time(1772323200), want: `<tg-time unix="1772323200">date</tg-time>`},
		{name: "time format", got: EscapeHTML("date").TimeFormat(1772323200, `MMM " yyyy`), want: `<tg-time unix="1772323200" format="MMM &quot; yyyy">date</tg-time>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("formatted HTML link = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
