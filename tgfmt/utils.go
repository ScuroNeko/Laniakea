package tgfmt

import "strings"

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}
func escapeRich(s string) Rich { return Rich(escapeHTML(s)) }

// EscapePunctuation escapes '.', '!' and '-' for MarkdownV2 fragments.
func EscapePunctuation(s string) string {
	symbols := []string{".", "!", "-"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, "\\"+symbol)
	}
	return s
}
