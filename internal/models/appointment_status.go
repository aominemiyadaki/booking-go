package models

import "time"

// AppointmentStatus — статус посещения (ожидает, подтверждено, отменено и т.д.).
//
// В старой схеме статус хранился как строка (VARCHAR) прямо в таблице appointments.
// В новой схеме статусы вынесены в отдельную справочную таблицу.
// Это позволяет:
//   - добавлять новые статусы без миграции
//   - переводить названия статусов на разные языки
//   - хранить дополнительные свойства статуса (цвет, иконка)
type AppointmentStatus struct {
	ID         string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Name       string    `json:"name" gorm:"uniqueIndex;size:255;not null"`
	Version    int64     `json:"version" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	Visibility bool      `json:"visibility" gorm:"default:true"`
}

func (AppointmentStatus) TableName() string {
	return "booking_bubuka.appointment_statuses"
}
