package logger

import (
	"io"
	"log/slog"
	"os"
)

func SetDefaultLogger(logLevelStr string) {
	var logLevel slog.Level
	var f io.Writer
	switch logLevelStr {
	case "debug":
		logLevel = slog.LevelDebug
		f = os.Stderr
	case "info":
		logLevel = slog.LevelInfo
		f = os.Stderr
	case "warn":
		logLevel = slog.LevelWarn
		f = os.Stderr
	case "error":
		logLevel = slog.LevelError
		f = os.Stderr
	default:
		logLevel = slog.LevelInfo
		f = os.Stderr
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: logLevel})))
}
