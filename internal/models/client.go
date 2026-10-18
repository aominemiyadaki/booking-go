package models

import "time"

// Client — клиент мастера.
//
// Изменения:
//   - phone + phone_normalized → phone_number (одно поле)
//   - добавлены birth_date, telegram_id, telegram_username
//   - добавлены version, visibility
type Client struct {
	ID               string     `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	MasterID         string     `json:"master_id" gorm:"not null"`
	Name             string     `json:"name" gorm:"size:255;not null"`
	PhoneNumber      string     `json:"phone_number" gorm:"uniqueIndex;size:20;not null"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	LastVisitAt      *time.Time `json:"last_visited,omitempty"`
	TotalVisits      int        `json:"total_visits" gorm:"default:1"`
	Notes            *string    `json:"notes,omitempty" gorm:"type:text"`
	IsBlocked        bool       `json:"is_blocked" gorm:"default:false"`
	Version          int64      `json:"version" gorm:"default:1"`
	Visibility       bool       `json:"visibility" gorm:"default:true"`
	BirthDate        *time.Time `json:"birth_date,omitempty" gorm:"type:date"`
	TelegramID       *int64     `json:"telegram_id,omitempty"`
	TelegramUsername *string    `json:"telegram_username,omitempty" gorm:"size:255"`
}

func (Client) TableName() string {
	return "booking_bubuka.clients"
}

// CreateClientRequest — ручное добавление клиента мастером.
type CreateClientRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Notes       *string `json:"notes,omitempty"`
}

// UpdateClientRequest — обновление данных клиента.
type UpdateClientRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Notes       *string `json:"notes,omitempty"`
}

// UpdateNotesRequest — обновление заметок о клиенте.
type UpdateNotesRequest struct {
	Notes string `json:"notes"`
}

// ToggleBlockRequest — блокировка/разблокировка клиента.
type ToggleBlockRequest struct {
	Block bool `json:"block"`
}

// FindOrCreateClientRequest — поиск или создание клиента (для клиентского флоу).
type FindOrCreateClientRequest struct {
	MasterID    string `json:"master_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// FindOrCreateClientResponse — ответ с ID клиента.
type FindOrCreateClientResponse struct {
	ClientID string `json:"client_id"`
}
