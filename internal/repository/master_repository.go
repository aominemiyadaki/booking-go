package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"booking-backend/internal/models"
)

type MasterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{db: db}
}

func (r *MasterRepository) Create(ctx context.Context, master *models.Master) error {
	return r.db.WithContext(ctx).Create(master).Error
}

func (r *MasterRepository) GetByEmail(ctx context.Context, email string) (*models.Master, error) {
	var master models.Master
	err := r.db.WithContext(ctx).
		Where("email = ? AND visibility = true", email).
		First(&master).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("мастер не найден")
		}
		return nil, err
	}

	if master.SpecializationID != "" {
		_ = r.db.WithContext(ctx).
			Table("booking_bubuka.specializations").
			Select("name").
			Where("id = ?", master.SpecializationID).
			Scan(&master.SpecializationName).Error
	}

	return &master, nil
}

func (r *MasterRepository) GetByID(ctx context.Context, id string) (*models.Master, error) {
	var master models.Master
	err := r.db.WithContext(ctx).
		Where("id = ? AND visibility = true", id).
		First(&master).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("мастер не найден")
		}
		return nil, err
	}

	if master.SpecializationID != "" {
		_ = r.db.WithContext(ctx).
			Table("booking_bubuka.specializations").
			Select("name").
			Where("id = ?", master.SpecializationID).
			Scan(&master.SpecializationName).Error
	}

	return &master, nil
}

func (r *MasterRepository) GetAll(ctx context.Context) ([]models.MasterPublic, error) {
	var masters []models.MasterPublic
	err := r.db.WithContext(ctx).
		Table("booking_bubuka.masters m").
		Select("m.id, m.full_name, m.specialization_id, s.name as specialization, m.avatar_url, m.phone_number").
		Joins("left join booking_bubuka.specializations s on m.specialization_id = s.id").
		Where("m.visibility = true").
		Order("m.full_name").
		Scan(&masters).Error
	return masters, err
}

// Update сохраняет данные профиля мастера.
func (r *MasterRepository) Update(ctx context.Context, m *models.Master) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MasterRepository) UpdateSessionToken(ctx context.Context, masterID, token string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Master{}).
		Where("id = ?", masterID).
		Update("session_token", token)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("мастер не найден")
	}
	return nil
}
