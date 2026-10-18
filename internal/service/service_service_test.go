package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-backend/internal/models"
	"booking-backend/internal/repository/mocks"
)

func TestServiceService_GetByMaster(t *testing.T) {
	repo := new(mocks.MockServiceRepository)
	svc := NewServiceService(repo)

	expected := []models.Service{{ID: "s1", Title: "Стрижка"}}
	repo.On("GetByMasterID", mock.Anything, "master-1").Return(expected, nil)

	services, err := svc.GetByMaster(context.Background(), "master-1")

	assert.NoError(t, err)
	assert.Len(t, services, 1)
	assert.Equal(t, "Стрижка", services[0].Title)
}

func TestServiceService_Create(t *testing.T) {
	repo := new(mocks.MockServiceRepository)
	svc := NewServiceService(repo)

	req := models.CreateServiceRequest{
		Title:           "Маникюр",
		DurationMinutes: 45,
		Price:           1500,
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Service")).Return(nil).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Service)
		s.ID = "s-new"
	})

	res, err := svc.Create(context.Background(), "master-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "s-new", res.ID)
	assert.Equal(t, "Маникюр", res.Title)
}

func TestServiceService_Update(t *testing.T) {
	repo := new(mocks.MockServiceRepository)
	svc := NewServiceService(repo)

	repo.On("Update", mock.Anything, mock.AnythingOfType("*models.Service")).Return(nil)

	err := svc.Update(context.Background(), "s1", "master-1", models.UpdateServiceRequest{Title: "Тест"})

	assert.NoError(t, err)
}

func TestServiceService_Delete(t *testing.T) {
	repo := new(mocks.MockServiceRepository)
	svc := NewServiceService(repo)

	repo.On("Delete", mock.Anything, "s1", "master-1").Return(nil)

	err := svc.Delete(context.Background(), "s1", "master-1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
