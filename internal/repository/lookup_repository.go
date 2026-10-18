package repository

import (
	"context"

	"gorm.io/gorm"

	"booking-backend/internal/models"
)

type LookupRepository struct {
	db *gorm.DB
}

func NewLookupRepository(db *gorm.DB) *LookupRepository {
	return &LookupRepository{db: db}
}

// Specializations
func (r *LookupRepository) GetSpecializations(ctx context.Context) ([]models.Specialization, error) {
	var list []models.Specialization
	err := r.db.WithContext(ctx).
		Where("visibility = true").
		Order("name").
		Find(&list).Error
	return list, err
}

func (r *LookupRepository) CreateSpecialization(ctx context.Context, s *models.Specialization) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *LookupRepository) UpdateSpecialization(ctx context.Context, s *models.Specialization) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *LookupRepository) DeleteSpecialization(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Specialization{}).
		Where("id = ?", id).
		Update("visibility", false).Error
}

// Currencies
func (r *LookupRepository) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	var list []models.Currency
	err := r.db.WithContext(ctx).
		Order("iso_code").
		Find(&list).Error
	return list, err
}

func (r *LookupRepository) CreateCurrency(ctx context.Context, c *models.Currency) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *LookupRepository) UpdateCurrency(ctx context.Context, c *models.Currency) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *LookupRepository) DeleteCurrency(ctx context.Context, id string) error {
	// В таблице currencies нет поля visibility, поэтому делаем физическое удаление (hard delete)
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Currency{}).Error
}

// AppointmentStatuses
func (r *LookupRepository) GetAppointmentStatuses(ctx context.Context) ([]models.AppointmentStatus, error) {
	var list []models.AppointmentStatus
	err := r.db.WithContext(ctx).
		Where("visibility = true").
		Order("name").
		Find(&list).Error
	return list, err
}

func (r *LookupRepository) CreateAppointmentStatus(ctx context.Context, s *models.AppointmentStatus) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *LookupRepository) UpdateAppointmentStatus(ctx context.Context, s *models.AppointmentStatus) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *LookupRepository) DeleteAppointmentStatus(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.AppointmentStatus{}).
		Where("id = ?", id).
		Update("visibility", false).Error
}

// BreakTypes
func (r *LookupRepository) GetBreakTypes(ctx context.Context) ([]models.BreakType, error) {
	var list []models.BreakType
	err := r.db.WithContext(ctx).
		Where("visibility = true").
		Order("name").
		Find(&list).Error
	return list, err
}

func (r *LookupRepository) CreateBreakType(ctx context.Context, b *models.BreakType) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *LookupRepository) UpdateBreakType(ctx context.Context, b *models.BreakType) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *LookupRepository) DeleteBreakType(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.BreakType{}).
		Where("id = ?", id).
		Update("visibility", false).Error
}
