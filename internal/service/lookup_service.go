package service

import (
	"context"

	"booking-backend/internal/models"
	"booking-backend/internal/repository"
)

type LookupService struct {
	repo repository.LookupRepositoryInterface
}

func NewLookupService(repo repository.LookupRepositoryInterface) *LookupService {
	return &LookupService{repo: repo}
}

// Specializations
func (s *LookupService) GetSpecializations(ctx context.Context) ([]models.Specialization, error) {
	return s.repo.GetSpecializations(ctx)
}

func (s *LookupService) CreateSpecialization(ctx context.Context, name string) (*models.Specialization, error) {
	item := &models.Specialization{Name: name}
	err := s.repo.CreateSpecialization(ctx, item)
	return item, err
}

func (s *LookupService) UpdateSpecialization(ctx context.Context, id, name string) error {
	item := &models.Specialization{ID: id, Name: name}
	return s.repo.UpdateSpecialization(ctx, item)
}

func (s *LookupService) DeleteSpecialization(ctx context.Context, id string) error {
	return s.repo.DeleteSpecialization(ctx, id)
}

// Currencies
func (s *LookupService) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	return s.repo.GetCurrencies(ctx)
}

func (s *LookupService) CreateCurrency(ctx context.Context, name, isoCode string) (*models.Currency, error) {
	item := &models.Currency{Name: name, ISOCode: isoCode}
	err := s.repo.CreateCurrency(ctx, item)
	return item, err
}

func (s *LookupService) UpdateCurrency(ctx context.Context, id, name, isoCode string) error {
	item := &models.Currency{ID: id, Name: name, ISOCode: isoCode}
	return s.repo.UpdateCurrency(ctx, item)
}

func (s *LookupService) DeleteCurrency(ctx context.Context, id string) error {
	return s.repo.DeleteCurrency(ctx, id)
}

// AppointmentStatuses
func (s *LookupService) GetAppointmentStatuses(ctx context.Context) ([]models.AppointmentStatus, error) {
	return s.repo.GetAppointmentStatuses(ctx)
}

func (s *LookupService) CreateAppointmentStatus(ctx context.Context, name string) (*models.AppointmentStatus, error) {
	item := &models.AppointmentStatus{Name: name}
	err := s.repo.CreateAppointmentStatus(ctx, item)
	return item, err
}

func (s *LookupService) UpdateAppointmentStatus(ctx context.Context, id, name string) error {
	item := &models.AppointmentStatus{ID: id, Name: name}
	return s.repo.UpdateAppointmentStatus(ctx, item)
}

func (s *LookupService) DeleteAppointmentStatus(ctx context.Context, id string) error {
	return s.repo.DeleteAppointmentStatus(ctx, id)
}

// BreakTypes
func (s *LookupService) GetBreakTypes(ctx context.Context) ([]models.BreakType, error) {
	return s.repo.GetBreakTypes(ctx)
}

func (s *LookupService) CreateBreakType(ctx context.Context, name string) (*models.BreakType, error) {
	item := &models.BreakType{Name: name}
	err := s.repo.CreateBreakType(ctx, item)
	return item, err
}

func (s *LookupService) UpdateBreakType(ctx context.Context, id, name string) error {
	item := &models.BreakType{ID: id, Name: name}
	return s.repo.UpdateBreakType(ctx, item)
}

func (s *LookupService) DeleteBreakType(ctx context.Context, id string) error {
	return s.repo.DeleteBreakType(ctx, id)
}
