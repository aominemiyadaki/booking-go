package models

import "time"

// BlacklistedClient — запись в чёрном списке.
//
// Изменения: упрощена до master_id + client_id.
// Старые поля (guest_id, reason, device_id, ip_address и т.д.) убраны.
type BlacklistedClient struct {
	ID         string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	MasterID   string    `json:"master_id" gorm:"not null"`
	ClientID   string    `json:"client_id" gorm:"not null"`
	Version    int64     `json:"version" gorm:"default:1"`
	Visibility bool      `json:"visibility" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (BlacklistedClient) TableName() string {
	return "booking_bubuka.blacklist"
}

// BlockClientRequest — запрос на блокировку клиента.
type BlockClientRequest struct {
	ClientID string `json:"client_id" binding:"required"`
}
