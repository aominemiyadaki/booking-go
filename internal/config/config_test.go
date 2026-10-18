package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	os.Setenv("DB_HOST", "customhost")
	os.Setenv("DB_PORT", "5433")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
	}()

	cfg := Load()

	assert.Equal(t, "customhost", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Contains(t, cfg.DSN(), "postgres://postgres:postgres@customhost:5433/booking?sslmode=disable")
}
