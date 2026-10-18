package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"booking-backend/internal/models"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(ctx context.Context, a *models.Appointment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Проверка на пересечение (overlapping)
		var count int64
		err := tx.Table("booking_bubuka.appointments").
			Joins("JOIN booking_bubuka.appointment_statuses ON booking_bubuka.appointments.status_id = booking_bubuka.appointment_statuses.id").
			Where("booking_bubuka.appointments.master_id = ? AND booking_bubuka.appointments.visibility = true AND booking_bubuka.appointments.start_time < ? AND booking_bubuka.appointments.end_time > ? AND booking_bubuka.appointment_statuses.name != 'отменен'",
				a.MasterID, a.EndTime, a.StartTime).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("выбранное время уже занято")
		}

		return tx.Create(a).Error
	})
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id string) (*models.Appointment, error) {
	var app models.Appointment
	err := r.db.WithContext(ctx).
		Where("id = ? AND visibility = true", id).
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("запись не найдена")
		}
		return nil, err
	}
	return &app, nil
}

func (r *AppointmentRepository) GetByMasterID(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error) {
	page.Normalize()

	var total int64
	err := r.db.WithContext(ctx).
		Table("booking_bubuka.appointments").
		Where("master_id = ? AND visibility = true", masterID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var appointments []models.AppointmentResponse
	err = r.db.WithContext(ctx).
		Table("booking_bubuka.appointments a").
		Select("a.id, a.master_id, a.client_id, c.name as client_name, a.service_id, s.title as service_name, a.start_time, a.end_time, a.status_id, ast.name as status_name, a.comment, a.created_at").
		Joins("LEFT JOIN booking_bubuka.clients c ON a.client_id = c.id").
		Joins("LEFT JOIN booking_bubuka.services s ON a.service_id = s.id").
		Joins("LEFT JOIN booking_bubuka.appointment_statuses ast ON a.status_id = ast.id").
		Where("a.master_id = ? AND a.visibility = true", masterID).
		Order("a.start_time DESC").
		Offset(page.Offset()).
		Limit(page.PageSize).
		Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (r *AppointmentRepository) GetByClientID(ctx context.Context, clientID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error) {
	page.Normalize()

	var total int64
	err := r.db.WithContext(ctx).
		Table("booking_bubuka.appointments").
		Where("client_id = ? AND visibility = true", clientID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var appointments []models.AppointmentResponse
	err = r.db.WithContext(ctx).
		Table("booking_bubuka.appointments a").
		Select("a.id, a.master_id, a.client_id, c.name as client_name, a.service_id, s.title as service_name, a.start_time, a.end_time, a.status_id, ast.name as status_name, a.comment, a.created_at").
		Joins("LEFT JOIN booking_bubuka.clients c ON a.client_id = c.id").
		Joins("LEFT JOIN booking_bubuka.services s ON a.service_id = s.id").
		Joins("LEFT JOIN booking_bubuka.appointment_statuses ast ON a.status_id = ast.id").
		Where("a.client_id = ? AND a.visibility = true", clientID).
		Order("a.start_time DESC").
		Offset(page.Offset()).
		Limit(page.PageSize).
		Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (r *AppointmentRepository) GetByDate(ctx context.Context, masterID string, date time.Time) ([]models.AppointmentResponse, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).UTC()
	endOfDay := startOfDay.AddDate(0, 0, 1)

	var appointments []models.AppointmentResponse
	err := r.db.WithContext(ctx).
		Table("booking_bubuka.appointments a").
		Select("a.id, a.master_id, a.client_id, c.name as client_name, a.service_id, s.title as service_name, a.start_time, a.end_time, a.status_id, ast.name as status_name, a.comment, a.created_at").
		Joins("LEFT JOIN booking_bubuka.clients c ON a.client_id = c.id").
		Joins("LEFT JOIN booking_bubuka.services s ON a.service_id = s.id").
		Joins("LEFT JOIN booking_bubuka.appointment_statuses ast ON a.status_id = ast.id").
		Where("a.master_id = ? AND a.start_time >= ? AND a.start_time < ? AND a.visibility = true", masterID, startOfDay, endOfDay).
		Order("a.start_time").
		Find(&appointments).Error
	return appointments, err
}

// Update сохраняет все изменения в записи.
func (r *AppointmentRepository) Update(ctx context.Context, a *models.Appointment) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *AppointmentRepository) UpdateStatus(ctx context.Context, id, masterID, statusID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("id = ? AND master_id = ?", id, masterID).
		Update("status_id", statusID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("запись не найдена или не принадлежит мастеру")
	}
	return nil
}

func (r *AppointmentRepository) Delete(ctx context.Context, id, masterID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Appointment{}).
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

func (r *AppointmentRepository) GetAppointmentResponseByID(ctx context.Context, id string) (*models.AppointmentResponse, error) {
	var app models.AppointmentResponse
	err := r.db.WithContext(ctx).
		Table("booking_bubuka.appointments a").
		Select("a.id, a.master_id, a.client_id, c.name as client_name, a.service_id, s.title as service_name, a.start_time, a.end_time, a.status_id, ast.name as status_name, a.comment, a.created_at").
		Joins("LEFT JOIN booking_bubuka.clients c ON a.client_id = c.id").
		Joins("LEFT JOIN booking_bubuka.services s ON a.service_id = s.id").
		Joins("LEFT JOIN booking_bubuka.appointment_statuses ast ON a.status_id = ast.id").
		Where("a.id = ? AND a.visibility = true", id).
		Take(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("запись не найдена")
		}
		return nil, err
	}
	return &app, nil
}
