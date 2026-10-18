package repository

import (
	"context"
	"errors"
	"math"

	"gorm.io/gorm"

	"booking-backend/internal/models"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) GetByMasterID(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.Client, int64, error) {
	page.Normalize()

	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.Client{}).
		Where("master_id = ? AND visibility = true", masterID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var clients []models.Client
	err = r.db.WithContext(ctx).
		Where("master_id = ? AND visibility = true", masterID).
		Order("name").
		Offset(page.Offset()).
		Limit(page.PageSize).
		Find(&clients).Error
	if err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (*models.Client, error) {
	var client models.Client
	err := r.db.WithContext(ctx).
		Where("id = ? AND visibility = true", id).
		First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("клиент не найден")
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) GetByIDAndMaster(ctx context.Context, id, masterID string) (*models.Client, error) {
	var client models.Client
	err := r.db.WithContext(ctx).
		Where("id = ? AND master_id = ? AND visibility = true", id, masterID).
		First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("клиент не найден")
		}
		return nil, err
	}
	return &client, nil
}

// GetByPhone находит клиента по номеру телефона.
func (r *ClientRepository) GetByPhone(ctx context.Context, phone string) (*models.Client, error) {
	var client models.Client
	err := r.db.WithContext(ctx).
		Where("phone_number = ? AND visibility = true", phone).
		First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("клиент не найден")
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) Create(ctx context.Context, c *models.Client) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ClientRepository) Update(ctx context.Context, c *models.Client) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *ClientRepository) Delete(ctx context.Context, id, masterID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Client{}).
		Where("id = ? AND master_id = ?", id, masterID).
		Update("visibility", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("клиент не найден или не принадлежит мастеру")
	}
	return nil
}

func (r *ClientRepository) UpdateNotes(ctx context.Context, id, masterID, notes string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Client{}).
		Where("id = ? AND master_id = ?", id, masterID).
		Update("notes", notes)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("клиент не найден или не принадлежит мастеру")
	}
	return nil
}

func (r *ClientRepository) ToggleBlock(ctx context.Context, id, masterID string, block bool) error {
	result := r.db.WithContext(ctx).
		Model(&models.Client{}).
		Where("id = ? AND master_id = ?", id, masterID).
		Update("is_blocked", block)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("клиент не найден или не принадлежит мастеру")
	}
	return nil
}

func (r *ClientRepository) IncrementVisits(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Client{}).
		Where("id = ?", id).
		UpdateColumn("total_visits", gorm.Expr("total_visits + ?", 1)).
		UpdateColumn("last_visit_at", gorm.Expr("NOW()")).Error
}

var _ = math.Ceil
