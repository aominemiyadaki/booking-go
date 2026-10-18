package domain

import (
	"errors"
	"time"
)

// Единые доменные константы приложения
const (
	// Статусы записей
	StatusCancelled   = "отменен"
	StatusConfirmed   = "подтвержден"
	StatusPending     = "ожидает"
	StatusCompleted   = "завершен"

	// Значения по умолчанию для слотов и времени
	DefaultTimezone        = "Europe/Moscow"
	DefaultSlotDuration    = 30 * time.Minute
	DefaultSlotStep        = 30 * time.Minute
	DefaultJWTExpireHours  = 24
	MinPasswordLength      = 8
	MinJWTSecretLength     = 32
	DefaultMaxBodySize     = 1 << 20 // 1 MB
	DefaultRateLimitPerMin = 10
)

// Единые доменные ошибки
var (
	ErrNotFound          = errors.New("запись не найдена")
	ErrAlreadyExists     = errors.New("запись уже существует")
	ErrUnauthorized      = errors.New("требуется авторизация")
	ErrForbidden         = errors.New("доступ запрещён")
	ErrInvalidInput      = errors.New("невалидный ввод")
	ErrSlotTaken         = errors.New("выбранное время уже занято")
	ErrWeakPassword      = errors.New("пароль не соответствует требованиям безопасности")
	ErrWeakJWTSecret     = errors.New("JWT_SECRET слишком короткий или ненадёжный")
	ErrInternalServer    = errors.New("внутренняя ошибка сервера")
)
