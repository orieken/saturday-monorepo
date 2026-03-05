package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	if logger == nil {
		t.Fatal("Expected logger to be created, got nil")
	}

	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "key") {
		t.Errorf("Expected log output to contain 'key', got: %s", output)
	}
}

func TestNewLoggerWithLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithLevel(&buf, slog.LevelDebug)

	if logger == nil {
		t.Fatal("Expected logger to be created, got nil")
	}

	logger.Debug("debug message")
	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected debug log output, got: %s", output)
	}
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		logFunc  func(*Logger)
		expected bool
	}{
		{
			name:  "Info level logs info",
			level: slog.LevelInfo,
			logFunc: func(l *Logger) {
				l.Info("info message")
			},
			expected: true,
		},
		{
			name:  "Info level skips debug",
			level: slog.LevelInfo,
			logFunc: func(l *Logger) {
				l.Debug("debug message")
			},
			expected: false,
		},
		{
			name:  "Error level logs error",
			level: slog.LevelError,
			logFunc: func(l *Logger) {
				l.Error("error message")
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLoggerWithLevel(&buf, tt.level)
			tt.logFunc(logger)

			output := buf.String()
			hasOutput := len(output) > 0

			if hasOutput != tt.expected {
				t.Errorf("Expected output=%v, got output=%v (len=%d)", tt.expected, hasOutput, len(output))
			}
		})
	}
}
