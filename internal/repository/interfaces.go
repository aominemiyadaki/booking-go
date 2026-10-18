// Package repository определяет интерфейсы для работы с базой данных.
//
// Интерфейс — это контракт: он описывает КАКИЕ методы есть,
// но не описывает КАК они реализованы.
//
// Зачем это нужно:
//  1. Сервисный слой зависит от интерфейса, а не от конкретной реализации.
//     Это значит, что вы можете заменить GORM на сырой SQL или MongoDB,
//     не меняя ни одной строки в сервисах.
//  2. В тестах можно подставить «мок» (подделку) вместо реальной базы данных.
package repository

import (
	"context"
	"time"

	"booking-backend/internal/models"
)

// MasterRepositoryInterface — интерфейс для работы с мастерами.
type MasterRepositoryInterface interface {
	Create(ctx context.Context, master *models.Master) error
	GetByEmail(ctx context.Context, email string) (*models.Master, error)
	GetByID(ctx context.Context, id string) (*models.Master, error)
	GetAll(ctx context.Context) ([]models.MasterPublic, error)
	Update(ctx context.Context, m *models.Master) error
	UpdateSessionToken(ctx context.Context, masterID, token string) error
}

// ServiceRepositoryInterface — интерфейс для работы с услугами.
type ServiceRepositoryInterface interface {
	GetByMasterID(ctx context.Context, masterID string) ([]models.Service, error)
	GetByID(ctx context.Context, id string) (*models.Service, error)
	Create(ctx context.Context, s *models.Service) error
	Update(ctx context.Context, s *models.Service) error
	Delete(ctx context.Context, id, masterID string) error
}

// ClientRepositoryInterface — интерфейс для работы с клиентами.
type ClientRepositoryInterface interface {
	GetByMasterID(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.Client, int64, error)
	GetByID(ctx context.Context, id string) (*models.Client, error)
	GetByIDAndMaster(ctx context.Context, id, masterID string) (*models.Client, error)
	GetByPhone(ctx context.Context, phone string) (*models.Client, error)
	Create(ctx context.Context, c *models.Client) error
	Update(ctx context.Context, c *models.Client) error
	Delete(ctx context.Context, id, masterID string) error
	UpdateNotes(ctx context.Context, id, masterID, notes string) error
	ToggleBlock(ctx context.Context, id, masterID string, block bool) error
	IncrementVisits(ctx context.Context, id string) error
}

// AppointmentRepositoryInterface — интерфейс для работы с записями.
type AppointmentRepositoryInterface interface {
	Create(ctx context.Context, a *models.Appointment) error
	GetByID(ctx context.Context, id string) (*models.Appointment, error)
	GetByMasterID(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error)
	GetByClientID(ctx context.Context, clientID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error)
	GetByDate(ctx context.Context, masterID string, date time.Time) ([]models.AppointmentResponse, error)
	GetAppointmentResponseByID(ctx context.Context, id string) (*models.AppointmentResponse, error)
	Update(ctx context.Context, a *models.Appointment) error
	UpdateStatus(ctx context.Context, id, masterID, statusID string) error
	Delete(ctx context.Context, id, masterID string) error
}

// BlacklistRepositoryInterface — интерфейс для работы с чёрным списком.
type BlacklistRepositoryInterface interface {
	GetByMasterID(ctx context.Context, masterID string) ([]models.BlacklistedClient, error)
	Add(ctx context.Context, b *models.BlacklistedClient) error
	Remove(ctx context.Context, id, masterID string) error
	IsBlocked(ctx context.Context, masterID, clientID string) (bool, error)
}

// ScheduleRepositoryInterface — интерфейс для работы с расписанием.
type ScheduleRepositoryInterface interface {
	GetWorkingSlotsForDate(ctx context.Context, masterID, date string) ([]models.TimeSlot, error)
	SaveTimeSlots(ctx context.Context, masterID string, dates []string, slots []models.TimeSlot) error
	DeleteWorkingDay(ctx context.Context, masterID, date string) error
	GetWorkingDaysForMonth(ctx context.Context, masterID, yearMonth string) ([]string, error)
	GetForMonth(ctx context.Context, masterID, yearMonth string) (*models.MonthlySchedule, error)

	// Перерывы мастера (Breaks)
	CreateBreak(ctx context.Context, b *models.MasterBreak) error
	GetBreaks(ctx context.Context, masterID string) ([]models.MasterBreak, error)
	DeleteBreak(ctx context.Context, id, masterID string) error
}

// LookupRepositoryInterface — интерфейс для справочников.
type LookupRepositoryInterface interface {
	GetSpecializations(ctx context.Context) ([]models.Specialization, error)
	CreateSpecialization(ctx context.Context, s *models.Specialization) error
	UpdateSpecialization(ctx context.Context, s *models.Specialization) error
	DeleteSpecialization(ctx context.Context, id string) error

	GetCurrencies(ctx context.Context) ([]models.Currency, error)
	CreateCurrency(ctx context.Context, c *models.Currency) error
	UpdateCurrency(ctx context.Context, c *models.Currency) error
	DeleteCurrency(ctx context.Context, id string) error

	GetAppointmentStatuses(ctx context.Context) ([]models.AppointmentStatus, error)
	CreateAppointmentStatus(ctx context.Context, s *models.AppointmentStatus) error
	UpdateAppointmentStatus(ctx context.Context, s *models.AppointmentStatus) error
	DeleteAppointmentStatus(ctx context.Context, id string) error

	GetBreakTypes(ctx context.Context) ([]models.BreakType, error)
	CreateBreakType(ctx context.Context, b *models.BreakType) error
	UpdateBreakType(ctx context.Context, b *models.BreakType) error
	DeleteBreakType(ctx context.Context, id string) error
}

