// Package logger builds the structured JSON logger.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a JSON logger writing to stdout at the given level.
func New(level string) *slog.Logger {
	l := slog.LevelInfo
	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}
