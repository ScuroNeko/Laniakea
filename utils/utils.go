package utils

import (
	"fmt"
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

func EscapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "_", `\_`)
	s = strings.ReplaceAll(s, "*", `\*`)
	s = strings.ReplaceAll(s, "[", `\[`)
	return strings.ReplaceAll(s, "`", "\\`")
}

func EscapeMarkdownV2(s string) string {
	symbols := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, fmt.Sprintf("\\%s", symbol))
	}
	return s
}
