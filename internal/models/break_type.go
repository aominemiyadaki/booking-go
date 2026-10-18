package models

import "time"

// BreakType — тип перерыва мастера (обед, личные дела, выходной и т.д.).
type BreakType struct {
	ID         string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Name       string    `json:"name" gorm:"uniqueIndex;size:255;not null"`
	Version    int64     `json:"version" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	Visibility bool      `json:"visibility" gorm:"default:true"`
}

func (BreakType) TableName() string {
	return "booking_bubuka.break_types"
}
