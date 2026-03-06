package utils

import (
	"os"
	"strings"

	"git.nix13.pw/scuroneko/slog"
)

func GetLoggerLevel() slog.LogLevel {
	level := slog.FATAL
	if os.Getenv("DEBUG") == "true" {
		level = slog.DEBUG
	}
	return level
}

// EscapeMarkdown Deprecated. Use MarkdownV2
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
	symbols := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!", "\\"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, "\\"+symbol)
	}
	return s
}
