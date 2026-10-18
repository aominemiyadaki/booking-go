package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtSecret := "test-secret-key-32-bytes-long!!"

	router := gin.New()
	router.Use(AuthMiddleware(jwtSecret))
	router.GET("/protected", func(c *gin.Context) {
		masterID, ok := getMasterID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"master_id": masterID})
	})

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Valid Token Header", func(t *testing.T) {
		claims := jwt.MapClaims{
			"master_id": "master-uuid-99",
			"exp":       time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(jwtSecret))

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "master-uuid-99")
	})

	t.Run("Valid Cookie Token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"master_id": "master-cookie-1",
			"exp":       time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(jwtSecret))

		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{Name: "master_jwt", Value: tokenString})
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "master-cookie-1")
	})
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeadersMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, "DENY", resp.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", resp.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "1; mode=block", resp.Header().Get("X-XSS-Protection"))
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware("https://example.com, https://app.example.com"))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("Allowed Origin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, "https://example.com", resp.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("OPTIONS Preflight", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	})
}

func TestRateLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	limiterMiddleware := RateLimiterMiddleware(2)
	router.Use(limiterMiddleware)
	router.GET("/limited", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 1st request -> OK
	req1 := httptest.NewRequest("GET", "/limited", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	assert.Equal(t, http.StatusOK, resp1.Code)

	// 2nd request -> OK
	req2 := httptest.NewRequest("GET", "/limited", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusOK, resp2.Code)

	// 3rd request -> Too Many Requests
	req3 := httptest.NewRequest("GET", "/limited", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	assert.Equal(t, http.StatusTooManyRequests, resp3.Code)
}
