package laniakea

import (
	"strings"

	"git.nix13.pw/scuroneko/laniakea/utils"
)

func Ptr[T any](v T) *T { return &v }
func Val[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

// EscapeMarkdown
// Deprecated. Use MarkdownV2
func EscapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "_", `\_`)
	s = strings.ReplaceAll(s, "*", `\*`)
	s = strings.ReplaceAll(s, "[", `\[`)
	return strings.ReplaceAll(s, "`", "\\`")
}

// EscapeHTML escapes special characters for Telegram HTML parse mode.
func EscapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// EscapeMarkdownV2 escapes special characters for Telegram MarkdownV2.
// https://core.telegram.org/bots/api#markdownv2-style
func EscapeMarkdownV2(s string) string {
	symbols := []string{"\\", "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, "\\"+symbol)
	}
	return s
}
func EscapePunctuation(s string) string {
	symbols := []string{".", "!", "-"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, "\\"+symbol)
	}
	return s
}

const (
	VersionString = utils.VersionString
	VersionMajor  = utils.VersionMajor
	VersionMinor  = utils.VersionMinor
	VersionPatch  = utils.VersionPatch
	VersionBeta   = utils.VersionBeta
)
