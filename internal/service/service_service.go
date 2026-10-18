package service

import (
	"context"

	"booking-backend/internal/models"
	"booking-backend/internal/repository"
)

type ServiceService struct {
	repo repository.ServiceRepositoryInterface
}

func NewServiceService(repo repository.ServiceRepositoryInterface) *ServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) GetByMaster(ctx context.Context, masterID string) ([]models.Service, error) {
	return s.repo.GetByMasterID(ctx, masterID)
}

func (s *ServiceService) GetByID(ctx context.Context, id string) (*models.Service, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ServiceService) Create(ctx context.Context, masterID string, req models.CreateServiceRequest) (*models.Service, error) {
	service := &models.Service{
		MasterID:        masterID,
		Title:           req.Title,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		Description:     req.Description,
		SortOrder:       req.SortOrder,
		CurrencyID:      req.CurrencyID,
	}

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceService) Update(ctx context.Context, id, masterID string, req models.UpdateServiceRequest) error {
	service := &models.Service{
		ID:              id,
		MasterID:        masterID,
		Title:           req.Title,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		Description:     req.Description,
		SortOrder:       req.SortOrder,
		CurrencyID:      req.CurrencyID,
	}

	return s.repo.Update(ctx, service)
}

func (s *ServiceService) Delete(ctx context.Context, id, masterID string) error {
	return s.repo.Delete(ctx, id, masterID)
}
