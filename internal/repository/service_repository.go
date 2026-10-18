package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"booking-backend/internal/models"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// GetByMasterID возвращает все видимые услуги мастера, отсортированные по sort_order.
func (r *ServiceRepository) GetByMasterID(ctx context.Context, masterID string) ([]models.Service, error) {
	var services []models.Service
	err := r.db.WithContext(ctx).Raw(`
		SELECT s.*, c.iso_code as currency_iso_code 
		FROM booking_bubuka.services s 
		LEFT JOIN booking_bubuka.currencies c ON s.currency_id = c.id 
		WHERE s.master_id = ? AND s.visibility = true 
		ORDER BY s.sort_order
	`, masterID).Scan(&services).Error
	return services, err
}

// GetByID находит услугу по ID.
func (r *ServiceRepository) GetByID(ctx context.Context, id string) (*models.Service, error) {
	var service models.Service
	result := r.db.WithContext(ctx).Raw(`
		SELECT s.*, c.iso_code as currency_iso_code 
		FROM booking_bubuka.services s 
		LEFT JOIN booking_bubuka.currencies c ON s.currency_id = c.id 
		WHERE s.id = ? AND s.visibility = true
	`, id).Scan(&service)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("услуга не найдена")
	}

	return &service, nil
}

// Create создаёт новую услугу.
func (r *ServiceRepository) Create(ctx context.Context, s *models.Service) error {
	return r.db.WithContext(ctx).Create(s).Error
}

// Update обновляет услугу.
//
// db.Save(s) — обновляет ВСЕ поля структуры в БД.
// GORM определяет, какую строку обновлять, по полю ID.
func (r *ServiceRepository) Update(ctx context.Context, s *models.Service) error {
	return r.db.WithContext(ctx).
		Model(s).
		Select("title", "price", "duration_minutes", "description", "sort_order", "currency_id").
		Updates(s).Error
}

// Delete выполняет «мягкое удаление» — ставит visibility = false.
// Данные остаются в базе, но перестают показываться в запросах.
func (r *ServiceRepository) Delete(ctx context.Context, id, masterID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Service{}).
		Where("id = ? AND master_id = ?", id, masterID).
		Update("visibility", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("услуга не найдена или не принадлежит мастеру")
	}
	return nil
}
