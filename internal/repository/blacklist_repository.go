package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"booking-backend/internal/models"
)

type BlacklistRepository struct {
	db *gorm.DB
}

func NewBlacklistRepository(db *gorm.DB) *BlacklistRepository {
	return &BlacklistRepository{db: db}
}

func (r *BlacklistRepository) GetByMasterID(ctx context.Context, masterID string) ([]models.BlacklistedClient, error) {
	var list []models.BlacklistedClient
	err := r.db.WithContext(ctx).
		Where("master_id = ? AND visibility = true", masterID).
		Find(&list).Error
	return list, err
}

func (r *BlacklistRepository) Add(ctx context.Context, b *models.BlacklistedClient) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *BlacklistRepository) Remove(ctx context.Context, id, masterID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.BlacklistedClient{}).
		Where("id = ? AND master_id = ?", id, masterID).
		Update("visibility", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("запись не найдена или не принадлежит мастеру")
	}
	return nil
}

func (r *BlacklistRepository) IsBlocked(ctx context.Context, masterID, clientID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.BlacklistedClient{}).
		Where("master_id = ? AND client_id = ? AND visibility = true", masterID, clientID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
