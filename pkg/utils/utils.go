// Contents: utility functions for the project
package utils

import (
	"log/slog"
	"strings"
)

// Lower converts a string to lowercase
func Lower(s string) string {
	return strings.ToLower(s)
}

// LevelStringToSlog converts a log level string to slog.Level(int)
func LevelStringToSlog(level string) slog.Level {
	switch Lower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "fatal":
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}
