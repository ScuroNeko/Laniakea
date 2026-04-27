package utils

import (
	"os"

	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

// GetLoggerLevel returns DEBUG when DEBUG=true in env, otherwise FATAL.
func GetLoggerLevel() sneklog.LogLevel {
	level := sneklog.FATAL
	if os.Getenv("DEBUG") == "true" {
		level = sneklog.DEBUG
	}
	return level
}

// CreateLogger creates a logger with the shared default policy:
// JSON stdout output, provided prefix, and provided level.
func CreateLogger(prefix string, level sneklog.LogLevel) *sneklog.Logger {
	logger := sneklog.CreateLogger().Level(level)
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
func CreateFileLogger(prefix string, level sneklog.LogLevel, filePath string) (*sneklog.Logger, error) {
	logger := CreateLogger(prefix, level)
	fileWriter, err := logger.CreateTextFileWriter(filePath)
	if err != nil {
		return logger, err
	}
	logger.AddWriter(fileWriter)
	return logger, nil
}
