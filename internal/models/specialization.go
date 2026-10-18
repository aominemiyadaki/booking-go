package models

import "time"

// Specialization — специализация мастера (парикмахер, маникюр и т.д.).
//
// Это справочная таблица: содержит фиксированный список специализаций.
// Мастер ссылается на неё через поле SpecializationID.
type Specialization struct {
	// ID — уникальный идентификатор (UUID).
	// gorm:"primaryKey" — говорит GORM, что это первичный ключ.
	// gorm:"default:gen_random_uuid()" — PostgreSQL сам сгенерирует UUID при INSERT.
	ID string `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`

	// Name — название специализации (например, "Парикмахер").
	// gorm:"uniqueIndex" — создаёт UNIQUE INDEX, чтобы не было дублей.
	Name string `json:"name" gorm:"uniqueIndex;size:255;not null"`

	// Version — версия записи для оптимистичной блокировки.
	// При каждом обновлении версия увеличивается на 1.
	// Если два запроса одновременно обновляют запись, тот, у которого
	// старая версия, получит ошибку — это защищает от потери данных.
	Version int64 `json:"version" gorm:"default:1"`

	// CreatedAt — дата создания. GORM заполняет автоматически.
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Visibility — «мягкое удаление». Вместо DELETE FROM мы ставим visibility=false.
	// Данные остаются в базе, но не показываются пользователю.
	Visibility bool `json:"visibility" gorm:"default:true"`
}

// TableName — указывает GORM имя таблицы в БД.
// Без этого метода GORM искал бы таблицу "specializations" в схеме public.
func (Specialization) TableName() string {
	return "booking_bubuka.specializations"
}
