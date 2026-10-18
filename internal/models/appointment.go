package models

import "time"

// Appointment — запись на приём (посещение).
//
// Изменения:
//   - status (text) → status_id (UUID, ссылка на appointment_statuses)
//   - убраны date, type, guest_token, technical_info, is_notified
//   - добавлены version, visibility
type Appointment struct {
	ID         string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	MasterID   string    `json:"master_id" gorm:"not null"`
	ClientID   string    `json:"client_id" gorm:"not null"`
	ServiceID  string    `json:"service_id" gorm:"not null"`
	StartTime  time.Time `json:"start_time" gorm:"not null"`
	EndTime    time.Time `json:"end_time" gorm:"not null"`
	StatusID   string    `json:"status_id" gorm:"not null"`
	Comment    *string   `json:"comment,omitempty" gorm:"type:text"`
	Version    int64     `json:"version" gorm:"default:1"`
	Visibility bool      `json:"visibility" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Appointment) TableName() string {
	return "booking_bubuka.appointments"
}

// BookAppointmentRequest — запрос на создание записи клиентом или мастером.
type BookAppointmentRequest struct {
	MasterID    string  `json:"master_id" binding:"required"`
	ClientID    string  `json:"client_id"`
	ClientName  string  `json:"client_name"`
	ClientPhone string  `json:"client_phone"`
	ServiceID string  `json:"service_id" binding:"required"`
	Date      string  `json:"date" binding:"required"`      // "YYYY-MM-DD"
	SlotTime  string  `json:"slot_time" binding:"required"` // "HH:MM"
	StatusID  string  `json:"status_id"`
	Comment   *string `json:"comment,omitempty"`
}

// UpdateStatusRequest — изменение статуса записи.
type UpdateStatusRequest struct {
	StatusID string `json:"status_id"`
	Status   string `json:"status"`
}

// AppointmentResponse — расширенный ответ для записи.
type AppointmentResponse struct {
	ID          string    `json:"id"`
	MasterID    string    `json:"master_id,omitempty"`
	ClientID    string    `json:"client_id,omitempty"`
	ClientName  string    `json:"client_name"`
	ServiceID   string    `json:"service_id,omitempty"`
	ServiceName string    `json:"service_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	StatusID    string    `json:"status_id,omitempty"`
	StatusName  string    `json:"status_name"`
	Comment     *string   `json:"comment,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	IsBlock     bool      `json:"is_block,omitempty"`
}
