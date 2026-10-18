package models

import "time"

// MasterBreak — перерыв мастера (конкретный интервал времени).
//
// В старой схеме перерывы были частью appointments (type = "block").
// В новой схеме перерывы — отдельная сущность со своим типом (BreakType).
type MasterBreak struct {
	ID          string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	MasterID    string    `json:"master_id" gorm:"not null"`
	BreakTypeID string    `json:"break_type_id" gorm:"not null"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	Version     int64     `json:"version" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	Visibility  bool      `json:"visibility" gorm:"default:true"`
}

func (MasterBreak) TableName() string {
	return "booking_bubuka.master_breaks"
}

// CreateMasterBreakRequest — запрос на создание перерыва.
type CreateMasterBreakRequest struct {
	BreakTypeID string    `json:"break_type_id" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}
