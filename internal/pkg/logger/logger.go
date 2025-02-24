package logger

import (
	"log/slog"
	"os"
)

func InitLogger(level slog.Level) *slog.Logger {
	options := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, options)

	return slog.New(handler)
}
