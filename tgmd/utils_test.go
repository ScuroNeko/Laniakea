package tgmd

import "testing"

func TestFormattingHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "bold", got: WithBold("text"), want: "*text*"},
		{name: "italic", got: WithItalic("text"), want: "_text_"},
		{name: "inline code", got: WithInlineCode("text"), want: "`text`"},
		{name: "link", got: WithLink("Laniakea", "https://example.test"), want: "[Laniakea](https://example.test)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("unexpected formatted text: got %q want %q", tc.got, tc.want)
			}
		})
	}
}
