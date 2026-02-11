package tgapi

import (
	"os"

	"git.nix13.pw/scuroneko/slog"
)

func GetLoggerLevel() slog.LogLevel {
	level := slog.FATAL
	if os.Getenv("DEBUG") == "true" {
		level = slog.DEBUG
	}
	return level
}
