// Package database предоставляет подключение к PostgreSQL через GORM.
//
// GORM — это ORM (Object-Relational Mapping) для Go.
// Вместо написания SQL-запросов руками, GORM позволяет работать
// с базой данных через Go-структуры.
//
// Пример: вместо "INSERT INTO masters (...) VALUES (...)"
// вы пишете db.Create(&master) — и GORM сам сгенерирует SQL.
package database

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"booking-backend/internal/config"
)

// New создаёт подключение к PostgreSQL через GORM.
//
// Что происходит внутри:
//  1. gorm.Open — открывает соединение с базой данных
//  2. postgres.Open(dsn) — говорит GORM использовать PostgreSQL
//  3. &gorm.Config{} — настройки GORM (логирование запросов и т.д.)
//
// Возвращает *gorm.DB — через этот объект делаются все запросы к БД.
func New(cfg *config.Config) (*gorm.DB, error) {
	// SECURITY: В production не логируем SQL-запросы (защита от утечки PII)
	logLevel := logger.Info
	if os.Getenv("GIN_MODE") == "release" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		// Logger: в dev-режиме показывает SQL, в production — только предупреждения и ошибки
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	// Проверяем, что база данных доступна.
	// db.DB() возвращает стандартный *sql.DB из пакета database/sql,
	// а Ping() отправляет простой запрос, чтобы убедиться, что БД отвечает.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("БД недоступна: %w", err)
	}

	slog.Info("Подключение к PostgreSQL через GORM установлено")
	return db, nil
}
