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

// CreateLogger creates a logger with the shared default policy:
// JSON stdout output, provided prefix, and provided level.
func CreateLogger(prefix string, level slog.LogLevel) *slog.Logger {
	logger := slog.CreateLogger().Level(level)
	if prefix != "" {
		logger.Prefix(prefix)
	}
	logger.AddWriter(logger.CreateJsonStdoutWriter())
	return logger
}

// CreateFileLogger creates a logger with the shared default policy and appends
// file output to the provided path.
//
// The returned logger is always non-nil. When file writer creation fails, the
// logger still writes to stdout and the error is returned to the caller.
func CreateFileLogger(prefix string, level slog.LogLevel, filePath string) (*slog.Logger, error) {
	logger := CreateLogger(prefix, level)
	fileWriter, err := logger.CreateTextFileWriter(filePath)
	if err != nil {
		return logger, err
	}
	logger.AddWriter(fileWriter)
	return logger, nil
}
