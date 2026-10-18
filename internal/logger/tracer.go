package logger

import (
	"context"
	"log/slog"
	"time"
)

// Trace логирует начало и конец выполнения функции, а также затраченное время.
// Использование: defer logger.Trace(ctx, "ИмяСервиса.ИмяМетода")()
func Trace(ctx context.Context, operation string, extraAttrs ...any) func() {
	start := time.Now()
	
	reqID := ""
	if val, ok := ctx.Value("request_id").(string); ok {
		reqID = val
	}

	attrs := []any{
		slog.String("operation", operation),
		slog.String("request_id", reqID),
	}
	attrs = append(attrs, extraAttrs...)

	slog.Info("Начало: "+operation, attrs...)

	return func() {
		duration := time.Since(start)
		endAttrs := append(attrs, slog.Duration("duration", duration), slog.Int64("duration_ms", duration.Milliseconds()))
		slog.Info("Завершение: "+operation, endAttrs...)
	}
}
