package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware проверяет JWT токен и устанавливает master_id в контекст.
// Поддерживает гибридное чтение токена: из HttpOnly cookie "master_jwt" или из заголовка Authorization: Bearer.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID, _ := c.Get("request_id")

		tokenString := ""
		if cookie, err := c.Cookie("master_jwt"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			slog.Warn("Попытка доступа без токена", "request_id", reqID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			slog.Warn("Предоставлен невалидный JWT токен", "error", err, "request_id", reqID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Warn("Невалидная структура Claims в токене", "request_id", reqID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидные claims"})
			c.Abort()
			return
		}

		masterID, ok := claims["master_id"].(string)
		if !ok {
			slog.Warn("В токене отсутствует master_id", "request_id", reqID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "отсутствует master_id в токене"})
			c.Abort()
			return
		}

		c.Set("master_id", masterID)
		c.Next()
	}
}

// SecurityHeadersMiddleware устанавливает заголовок безопасности OWASP для защиты от XSS, Clickjacking и MIME-sniffing.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

// getMasterID извлекает master_id из контекста (безопасно, без паники).
func getMasterID(c *gin.Context) (string, bool) {
	val, exists := c.Get("master_id")
	if !exists {
		return "", false
	}
	masterID, ok := val.(string)
	if !ok {
		return "", false
	}
	return masterID, true
}

// requireMasterID — helper для хендлеров: извлекает masterID или отдаёт 401.
func requireMasterID(c *gin.Context) (string, bool) {
	masterID, ok := getMasterID(c)
	if !ok {
		reqID, _ := c.Get("request_id")
		slog.Error("master_id отсутствует в контексте", "request_id", reqID)
		
		c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидная сессия"})
		c.Abort()
		return "", false
	}
	return masterID, true
}

// RequestIDMiddleware добавляет уникальный ID каждому запросу.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = c.GetHeader("X-Guest-Token")
		}
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// CORSMiddleware — настраивает заголовки CORS.
// В production разрешает только конкретные домены из конфигурации.
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Проверяем, разрешён ли Origin
		allowed := false
		for _, o := range origins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed && origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Guest-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// MaxBodySizeMiddleware ограничивает размер тела запроса для защиты от OOM-атак.
func MaxBodySizeMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}

type ipRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
	stopChan chan struct{}
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	limiter := &ipRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		stopChan: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				limiter.cleanup()
			case <-limiter.stopChan:
				return
			}
		}
	}()

	return limiter
}

func (l *ipRateLimiter) stop() {
	close(l.stopChan)
}

func (l *ipRateLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for ip, times := range l.requests {
		valid := make([]time.Time, 0, len(times))
		for _, t := range times {
			if now.Sub(t) < l.window {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(l.requests, ip)
		} else {
			l.requests[ip] = valid
		}
	}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	valid := make([]time.Time, 0, len(l.requests[ip]))
	for _, t := range l.requests[ip] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.limit {
		l.requests[ip] = valid
		return false
	}

	valid = append(valid, now)
	l.requests[ip] = valid
	return true
}

// RateLimiterMiddleware ограничивает количество запросов с одного IP (защита от брутфорса).
func RateLimiterMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiter := newIPRateLimiter(requestsPerMinute, time.Minute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip) {
			reqID, _ := c.Get("request_id")
			slog.Warn("Превышен лимит запросов (Rate Limit Exceeded)", "ip", ip, "request_id", reqID)

			c.Header("Retry-After", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "слишком много запросов, повторите попытку через минуту",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

