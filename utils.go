package laniakea

import (
	"strings"

	"git.nix13.pw/scuroneko/laniakea/utils"
)

// Ptr returns a pointer to v.
func Ptr[T any](v T) *T { return &v }

// Val returns dereferenced pointer value or def when p is nil.
func Val[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

// EscapeMarkdown escapes special characters for legacy Telegram Markdown.
// Deprecated: Use EscapeMarkdownV2.
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

// EscapePunctuation escapes '.', '!' and '-' for MarkdownV2 fragments.
func EscapePunctuation(s string) string {
	symbols := []string{".", "!", "-"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, "\\"+symbol)
	}
	return s
}

const (
	// VersionString re-exports the module version string.
	VersionString = utils.VersionString
	// VersionMajor re-exports the module major version.
	VersionMajor = utils.VersionMajor
	// VersionMinor re-exports the module minor version.
	VersionMinor = utils.VersionMinor
	// VersionPatch re-exports the module patch version.
	VersionPatch = utils.VersionPatch
	// VersionBeta re-exports the module prerelease counter.
	VersionBeta = utils.VersionBeta
)
