package models

import "time"

// WorkingDay — рабочий день мастера.
type WorkingDay struct {
	ID         string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	MasterID   string    `json:"master_id" gorm:"not null"`
	Date       time.Time `json:"date" gorm:"type:date;not null"`
	Version    int64     `json:"version" gorm:"default:1"`
	Visibility bool      `json:"visibility" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (WorkingDay) TableName() string {
	return "booking_bubuka.working_days"
}

// WorkingSlot — рабочий интервал дня (например, 09:00–18:00).
type WorkingSlot struct {
	ID           string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	WorkingDayID string    `json:"working_day_id" gorm:"not null"`
	StartTime    string    `json:"start_time" gorm:"type:time;not null"` // "HH:MM:SS"
	EndTime      string    `json:"end_time" gorm:"type:time;not null"`   // "HH:MM:SS"
	Version      int64     `json:"version" gorm:"default:1"`
	Visibility   bool      `json:"visibility" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (WorkingSlot) TableName() string {
	return "booking_bubuka.working_slots"
}

// TimeSlot — упрощённый интервал для API-ответов.
type TimeSlot struct {
	Start string `json:"start"` // "HH:MM"
	End   string `json:"end"`   // "HH:MM"
}

// MonthlySchedule — ответ API: расписание мастера на месяц.
type MonthlySchedule struct {
	MasterID    string                `json:"master_id"`
	YearMonth   string                `json:"year_month"`
	WorkingDays map[string][]TimeSlot `json:"working_days"`
}

// SaveTimeSlotsRequest — сохранение слотов для списка дат.
type SaveTimeSlotsRequest struct {
	Dates []string   `json:"dates" binding:"required,min=1"`
	Slots []TimeSlot `json:"slots" binding:"required,min=1"`
}

// AvailableSlotsQuery — параметры запроса доступных слотов.
type AvailableSlotsQuery struct {
	MasterID  string `form:"master_id" binding:"required"`
	Date      string `form:"date" binding:"required"` // "YYYY-MM-DD"
	ServiceID string `form:"service_id" binding:"omitempty"`
}

// WorkingDaysQuery — параметры запроса рабочих дней месяца.
type WorkingDaysQuery struct {
	MasterID string `form:"master_id" binding:"required"`
	Month    string `form:"month" binding:"required"` // "YYYY-MM"
}
