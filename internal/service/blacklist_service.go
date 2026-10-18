package service

import (
	"context"

	"booking-backend/internal/models"
	"booking-backend/internal/repository"
)

type BlacklistService struct {
	repo repository.BlacklistRepositoryInterface
}

func NewBlacklistService(repo repository.BlacklistRepositoryInterface) *BlacklistService {
	return &BlacklistService{repo: repo}
}

func (s *BlacklistService) GetAll(ctx context.Context, masterID string) ([]models.BlacklistedClient, error) {
	return s.repo.GetByMasterID(ctx, masterID)
}

func (s *BlacklistService) Block(ctx context.Context, masterID string, req models.BlockClientRequest) error {
	client := &models.BlacklistedClient{
		MasterID: masterID,
		ClientID: req.ClientID,
	}
	return s.repo.Add(ctx, client)
}

func (s *BlacklistService) Unblock(ctx context.Context, id, masterID string) error {
	return s.repo.Remove(ctx, id, masterID)
}

// Check проверяет, заблокирован ли пользователь по телефону или токену (заглушка для совместимости)
func (s *BlacklistService) Check(ctx context.Context, masterID, phone, guestToken string) (bool, error) {
	// В новой упрощенной схеме мы блокируем по client_id.
	// Для обратной совместимости возвращаем false.
	return false, nil
}
