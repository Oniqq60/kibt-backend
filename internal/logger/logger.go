package logger

import (
	"io"
	"log/slog"
	"strings"
)

func Init(level string, output io.Writer) {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: false,
	})

	slog.SetDefault(slog.New(handler))
}
