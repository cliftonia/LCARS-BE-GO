// Package logger provides structured logging configuration.
package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

func init() {
	// Create a default logger with JSON handler
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	defaultLogger = slog.New(handler)
}

// Get returns the default logger
func Get() *slog.Logger {
	return defaultLogger
}

// New creates a new logger with the specified environment
func New(environment string) *slog.Logger {
	var level slog.Level
	var handler slog.Handler

	// Set log level based on environment
	switch environment {
	case "development":
		level = slog.LevelDebug
		// Use text handler for development (more readable)
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	case "production":
		level = slog.LevelInfo
		// Use JSON handler for production (structured)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	default:
		level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	logger := slog.New(handler)
	defaultLogger = logger
	return logger
}
