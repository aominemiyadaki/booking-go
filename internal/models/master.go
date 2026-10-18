package models

import "time"

// Master — мастер (специалист) в системе.
//
// Основные изменения по сравнению со старой схемой:
//   - name → full_name (более точное имя)
//   - specialization (текст) → specialization_id (ссылка на таблицу specializations)
//   - добавлен phone_number
//   - добавлены version, visibility
type Master struct {
	ID               string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Version          int64     `json:"version" gorm:"default:1"`
	Email            string    `json:"email" gorm:"uniqueIndex;size:64;not null"`
	PasswordHash     string    `json:"-" gorm:"column:password_hash;not null"` // json:"-" — никогда не отдаём наружу
	FullName         string    `json:"full_name" gorm:"column:full_name;size:100;not null"`
	SpecializationID string    `json:"specialization_id" gorm:"not null"`
	AvatarURL        *string   `json:"avatar_url,omitempty"`
	SessionToken     *string   `json:"session_token,omitempty"`
	PhoneNumber      string    `json:"phone_number" gorm:"uniqueIndex;size:20;not null"`
	Visibility       bool      `json:"visibility" gorm:"default:true"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Virtual field (не хранится в БД)
	SpecializationName string `json:"specialization" gorm:"-"`
}

func (Master) TableName() string {
	return "booking_bubuka.masters"
}

// MasterPublic — публичная информация о мастере (для клиентов).
// Не содержит пароль, email, session_token.
type MasterPublic struct {
	ID             string  `json:"id"`
	FullName       string  `json:"full_name"`
	SpecializationID string  `json:"-"`
	Specialization string  `json:"specialization"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
	PhoneNumber    string  `json:"phone_number"`
}

// RegisterRequest — тело запроса на регистрацию мастера.
type RegisterRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	FullName         string `json:"full_name" binding:"required"`
	PhoneNumber      string `json:"phone_number" binding:"required"`
	SpecializationID string `json:"specialization_id" binding:"required"`
}

// LoginRequest — тело запроса на логин.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse — ответ с JWT токеном.
type AuthResponse struct {
	Token  string       `json:"token"`
	Master MasterPublic `json:"master"`
}

// UpdateSessionTokenRequest — обновление push-токена.
type UpdateSessionTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// UpdateMasterProfileRequest — запрос на обновление профиля мастера.
type UpdateMasterProfileRequest struct {
	ID               string `json:"id" binding:"required,uuid"`
	FullName         string `json:"full_name" binding:"required"`
	SpecializationID string `json:"specialization_id" binding:"required,uuid"`
	PhoneNumber      string `json:"phone_number" binding:"required"`
}
