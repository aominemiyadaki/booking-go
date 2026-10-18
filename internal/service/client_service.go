package service

import (
	"context"

	"booking-backend/internal/models"
	"booking-backend/internal/repository"
)

type ClientService struct {
	clientRepo      repository.ClientRepositoryInterface
	appointmentRepo repository.AppointmentRepositoryInterface
}

func NewClientService(clientRepo repository.ClientRepositoryInterface, appointmentRepo repository.AppointmentRepositoryInterface) *ClientService {
	return &ClientService{
		clientRepo:      clientRepo,
		appointmentRepo: appointmentRepo,
	}
}

func (s *ClientService) GetByMaster(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.Client, int64, error) {
	return s.clientRepo.GetByMasterID(ctx, masterID, page)
}

func (s *ClientService) GetByID(ctx context.Context, id, masterID string) (*models.Client, error) {
	return s.clientRepo.GetByIDAndMaster(ctx, id, masterID)
}

func (s *ClientService) Create(ctx context.Context, masterID string, req models.CreateClientRequest) (*models.Client, error) {
	c := &models.Client{
		MasterID:    masterID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Notes:       req.Notes,
	}

	if err := s.clientRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ClientService) Update(ctx context.Context, id, masterID string, req models.UpdateClientRequest) error {
	c := &models.Client{
		ID:          id,
		MasterID:    masterID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Notes:       req.Notes,
	}
	return s.clientRepo.Update(ctx, c)
}

func (s *ClientService) Delete(ctx context.Context, id, masterID string) error {
	return s.clientRepo.Delete(ctx, id, masterID)
}

func (s *ClientService) UpdateNotes(ctx context.Context, id, masterID, notes string) error {
	return s.clientRepo.UpdateNotes(ctx, id, masterID, notes)
}

func (s *ClientService) ToggleBlock(ctx context.Context, id, masterID string, block bool) error {
	return s.clientRepo.ToggleBlock(ctx, id, masterID, block)
}

// GetHistory возвращает историю визитов (записей) клиента.
func (s *ClientService) GetHistory(ctx context.Context, clientID, masterID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error) {
	// В репозитории записей у нас уже есть метод GetByClientID
	return s.appointmentRepo.GetByClientID(ctx, clientID, page)
}

func (s *ClientService) FindOrCreate(ctx context.Context, req models.FindOrCreateClientRequest) (*models.Client, error) {
	normalized := normalizePhone(req.PhoneNumber)
	client, err := s.clientRepo.GetByPhone(ctx, normalized)
	if err == nil {
		return client, nil
	}

	c := &models.Client{
		MasterID:    req.MasterID,
		Name:        req.Name,
		PhoneNumber: normalized,
	}

	if err := s.clientRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}
