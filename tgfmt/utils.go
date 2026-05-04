package tgfmt

import "strings"

// EscapePunctuation escapes '.', '!' and '-' for MarkdownV2 fragments.
func EscapePunctuation(s string) string {
	symbols := []string{".", "!", "-"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, "\\"+symbol)
	}
	return s
}
