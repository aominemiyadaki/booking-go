package models

// Currency — валюта для цен на услуги (рубль, тенге, доллар и т.д.).
//
// Это справочная таблица. Услуга (Service) ссылается на неё
// через поле CurrencyID, чтобы указать, в какой валюте цена.
type Currency struct {
	ID      string `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Name    string `json:"name" gorm:"uniqueIndex;size:255;not null"`
	ISOCode string `json:"iso_code" gorm:"uniqueIndex;size:3;not null"` // "KZT", "RUB", "USD"
	Version int64  `json:"version" gorm:"default:1"`
}

func (Currency) TableName() string {
	return "booking_bubuka.currencies"
}
