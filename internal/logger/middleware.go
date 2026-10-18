package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware — паттерн Middleware.
func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.String("client_ip", c.ClientIP()),
		}

		if reqID, ok := c.Get("request_id"); ok {
			attrs = append(attrs, slog.String("request_id", reqID.(string)))
		}

		if status >= 500 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
			log.LogAttrs(c.Request.Context(), slog.LevelError, "HTTP Request Failed", attrs...)
		} else {
			log.LogAttrs(c.Request.Context(), slog.LevelInfo, "HTTP Request Processed", attrs...)
		}
	}
}
