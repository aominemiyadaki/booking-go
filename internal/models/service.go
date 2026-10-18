package models

import "time"

// Service — услуга мастера (стрижка, маникюр и т.д.).
//
// Изменения по сравнению со старой схемой:
//   - добавлен currency_id (ссылка на таблицу валют)
//   - добавлены version, visibility
type Service struct {
	ID              string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	MasterID        string    `json:"master_id" gorm:"not null"`
	Title           string    `json:"title" gorm:"size:255;not null"`
	Price           float64   `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
	DurationMinutes int       `json:"duration_minutes" gorm:"not null;default:0"`
	Description     string    `json:"description" gorm:"type:text;default:''"`
	SortOrder       int       `json:"sort_order" gorm:"default:0"`
	CurrencyID      string    `json:"currency_id" gorm:"not null"`
	Version         int64     `json:"version" gorm:"default:1"`
	Visibility      bool      `json:"visibility" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Virtual field (не хранится в БД, только для чтения/сканирования)
	CurrencyISOCode string    `json:"currency_iso_code" gorm:"->"`
}

func (Service) TableName() string {
	return "booking_bubuka.services"
}

// CreateServiceRequest — добавление услуги.
type CreateServiceRequest struct {
	Title           string  `json:"title" binding:"required"`
	Price           float64 `json:"price" binding:"gte=0"`
	DurationMinutes int     `json:"duration_minutes" binding:"required,gt=0"`
	Description     string  `json:"description"`
	SortOrder       int     `json:"sort_order"`
	CurrencyID      string  `json:"currency_id" binding:"required"`
}

// UpdateServiceRequest — обновление услуги.
type UpdateServiceRequest struct {
	Title           string  `json:"title" binding:"required"`
	Price           float64 `json:"price" binding:"gte=0"`
	DurationMinutes int     `json:"duration_minutes" binding:"required,gt=0"`
	Description     string  `json:"description"`
	SortOrder       int     `json:"sort_order"`
	CurrencyID      string  `json:"currency_id" binding:"required"`
}
