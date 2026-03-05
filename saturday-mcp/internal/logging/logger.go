package logging

import (
	"io"
	"log/slog"
)

// Logger wraps slog for structured logging
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger
func NewLogger(w io.Writer) *Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithLevel creates a logger with a specific log level
func NewLoggerWithLevel(w io.Writer, level slog.Level) *Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	return &Logger{
		Logger: slog.New(handler),
	}
}
