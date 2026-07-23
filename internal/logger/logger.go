package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

func init() {
	level := getLogLevel()
	format := getLogFormat()

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
}

func getLogLevel() slog.Level {
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getLogFormat() string {
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	if format == "json" {
		return "json"
	}
	return "text"
}
