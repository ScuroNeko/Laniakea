package utils

import (
	"os"

	"git.nix13.pw/scuroneko/slog"
)

// GetLoggerLevel returns DEBUG when DEBUG=true in env, otherwise FATAL.
func GetLoggerLevel() slog.LogLevel {
	level := slog.FATAL
	if os.Getenv("DEBUG") == "true" {
		level = slog.DEBUG
	}
	return level
}
