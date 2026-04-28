package utils

import (
	"os"

	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
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
func CreateLogger(
	name string, level sneklog.LogLevel,
	format LogFormat, formatter *sneklog.Formatter,
) *sneklog.Logger {
	logger := sneklog.NewLogger().SetLevel(level)
	if name != "" {
		logger.SetName(name)
	}
	switch format {
	case LogFormatJSON:
		writer := logger.CreateJsonStdoutWriter()
		if formatter != nil {
			writer.SetFormatter(formatter)
		}
		logger.AddWriters(writer)
	default:
		writer := logger.CreateTextStdoutWriter()
		if formatter != nil {
			writer.SetFormatter(formatter)
		}
		logger.AddWriters(writer)
	}
	return logger
}

// CreateFileLogger creates a logger with the shared default policy and appends
// file output to the provided path.
//
// The returned logger is always non-nil. When file writer creation fails, the
// logger still writes to stdout and the error is returned to the caller.
func CreateFileLogger(
	prefix string, level sneklog.LogLevel, filePath string,
	format LogFormat, formatter *sneklog.Formatter,
) (*sneklog.Logger, error) {
	logger := CreateLogger(prefix, level, format, formatter)

	switch format {
	case LogFormatJSON:
		writer, err := logger.CreateJsonFileWriter(filePath)
		if err != nil {
			return logger, err
		}
		if formatter != nil {
			writer.SetFormatter(formatter)
		}
		logger.AddWriters(writer)
	default:
		writer, err := logger.CreateTextFileWriter(filePath)
		if err != nil {
			return logger, err
		}
		if formatter != nil {
			writer.SetFormatter(formatter)
		}
		logger.AddWriters(writer)
	}
	return logger, nil
}
