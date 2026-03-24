package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.nix13.pw/scuroneko/slog"
)

func TestCreateFileLoggerWritesToConfiguredFile(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "main.log")

	logger, err := CreateFileLogger("TEST", slog.DEBUG, logPath)
	if err != nil {
		t.Fatalf("CreateFileLogger returned error: %v", err)
	}
	logger.Infoln("hello from file logger")
	if err := logger.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if !strings.Contains(string(data), "hello from file logger") {
		t.Fatalf("expected log message in file, got %q", string(data))
	}
	if !strings.Contains(string(data), "[TEST]") {
		t.Fatalf("expected prefix in file, got %q", string(data))
	}
}
