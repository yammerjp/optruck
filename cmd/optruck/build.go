package optruck

import (
	"io"
	"log/slog"
	"os"
)

func (cli *CLI) buildLogger() *slog.Logger {
	var logLevel slog.Level
	var f io.Writer
	switch cli.LogLevel {
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
		f = io.Discard
	}

	return slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: logLevel}))
}
