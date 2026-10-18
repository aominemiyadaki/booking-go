// Package logger предоставляет настройки логирования.
package logger

import (
	"log/slog"
	"os"
)

// New создаёт логгер, который пишет всё в формате JSON в консоль (stdout).
func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler)
}
