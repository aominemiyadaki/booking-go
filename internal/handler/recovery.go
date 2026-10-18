package handler

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware — спасает сервер от падения (crash).
//
// В Go, если случается "panic" (например, деление на ноль или обращение
// по нулевому указателю), программа немедленно завершается.
// Это middleware использует "recover()", чтобы поймать панику, 
// записать её в Kibana вместе с request_id, и вернуть клиенту статус 500,
// а не просто разорвать соединение.
func RecoveryMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// defer выполняется в самом конце обработки запроса,
		// либо в момент возникновения паники.
		defer func() {
			if err := recover(); err != nil {
				// Достаём request_id для привязки лога
				reqID, _ := c.Get("request_id")
				
				// Логируем саму панику и стек-трейс (строки кода, где это случилось)
				log.Error("Критическая ошибка (Panic)!",
					"error", err,
					"request_id", reqID,
					"path", c.Request.URL.Path,
					"stack", string(debug.Stack()), // Stack trace поможет понять, где баг
				)

				// Отвечаем клиенту корректно
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "внутренняя ошибка сервера",
				})
			}
		}()

		// Идём дальше
		c.Next()
	}
}
