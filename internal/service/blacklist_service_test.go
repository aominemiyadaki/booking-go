package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-backend/internal/models"
	"booking-backend/internal/repository/mocks"
)

func TestBlacklistService(t *testing.T) {
	repo := new(mocks.MockBlacklistRepository)
	svc := NewBlacklistService(repo)

	ctx := context.Background()

	t.Run("GetAll", func(t *testing.T) {
		expected := []models.BlacklistedClient{{ID: "b1", ClientID: "c1"}}
		repo.On("GetByMasterID", ctx, "master-1").Return(expected, nil)

		list, err := svc.GetAll(ctx, "master-1")

		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("Block", func(t *testing.T) {
		req := models.BlockClientRequest{ClientID: "c1"}
		repo.On("Add", ctx, mock.AnythingOfType("*models.BlacklistedClient")).Return(nil)

		err := svc.Block(ctx, "master-1", req)

		assert.NoError(t, err)
	})

	t.Run("Unblock", func(t *testing.T) {
		repo.On("Remove", ctx, "b1", "master-1").Return(nil)
		err := svc.Unblock(ctx, "b1", "master-1")
		assert.NoError(t, err)
	})

	t.Run("Check", func(t *testing.T) {
		blocked, err := svc.Check(ctx, "master-1", "+79991112233", "token123")
		assert.NoError(t, err)
		assert.False(t, blocked)
	})
}
