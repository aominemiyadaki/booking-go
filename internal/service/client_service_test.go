package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-backend/internal/models"
	"booking-backend/internal/repository/mocks"
)

func TestClientService_GetByMaster(t *testing.T) {
	clientRepo := new(mocks.MockClientRepository)
	appRepo := new(mocks.MockAppointmentRepository)

	svc := NewClientService(clientRepo, appRepo)

	page := models.PaginationRequest{Page: 1, PageSize: 10}
	expectedClients := []models.Client{{ID: "c1", Name: "Клиент 1"}}

	clientRepo.On("GetByMasterID", mock.Anything, "master-1", page).Return(expectedClients, int64(1), nil)

	clients, total, err := svc.GetByMaster(context.Background(), "master-1", page)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, clients, 1)
	clientRepo.AssertExpectations(t)
}

func TestClientService_GetByID(t *testing.T) {
	clientRepo := new(mocks.MockClientRepository)
	appRepo := new(mocks.MockAppointmentRepository)

	svc := NewClientService(clientRepo, appRepo)

	clientRepo.On("GetByIDAndMaster", mock.Anything, "c1", "master-1").Return(&models.Client{ID: "c1", Name: "Клиент 1"}, nil)

	client, err := svc.GetByID(context.Background(), "c1", "master-1")

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "Клиент 1", client.Name)
}

func TestClientService_CreateAndUpdate(t *testing.T) {
	clientRepo := new(mocks.MockClientRepository)
	appRepo := new(mocks.MockAppointmentRepository)

	svc := NewClientService(clientRepo, appRepo)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		req := models.CreateClientRequest{Name: "Новый", PhoneNumber: "79998887766"}
		clientRepo.On("Create", ctx, mock.AnythingOfType("*models.Client")).Return(nil).Run(func(args mock.Arguments) {
			c := args.Get(1).(*models.Client)
			c.ID = "c-created"
		})

		c, err := svc.Create(ctx, "master-1", req)
		assert.NoError(t, err)
		assert.Equal(t, "c-created", c.ID)
	})

	t.Run("Update", func(t *testing.T) {
		req := models.UpdateClientRequest{Name: "Обновленный", PhoneNumber: "79998887766"}
		clientRepo.On("Update", ctx, mock.AnythingOfType("*models.Client")).Return(nil)

		err := svc.Update(ctx, "c1", "master-1", req)
		assert.NoError(t, err)
	})

	t.Run("GetHistory", func(t *testing.T) {
		page := models.PaginationRequest{Page: 1, PageSize: 10}
		appRepo.On("GetByClientID", ctx, "c1", page).Return([]models.AppointmentResponse{{ID: "a1"}}, int64(1), nil)

		apps, total, err := svc.GetHistory(ctx, "c1", "master-1", page)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, apps, 1)
	})
}

func TestClientService_FindOrCreate(t *testing.T) {
	clientRepo := new(mocks.MockClientRepository)
	appRepo := new(mocks.MockAppointmentRepository)

	svc := NewClientService(clientRepo, appRepo)

	t.Run("Existing Client", func(t *testing.T) {
		req := models.FindOrCreateClientRequest{PhoneNumber: "+79991234567", Name: "Тест"}
		clientRepo.On("GetByPhone", mock.Anything, "79991234567").Return(&models.Client{
			ID:          "c-existing",
			PhoneNumber: "79991234567",
			Name:        "Тест",
		}, nil)

		c, err := svc.FindOrCreate(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "c-existing", c.ID)
	})

	t.Run("New Client", func(t *testing.T) {
		req := models.FindOrCreateClientRequest{PhoneNumber: "+79990001122", Name: "Новый Клиент"}
		clientRepo.On("GetByPhone", mock.Anything, "79990001122").Return(nil, errors.New("not found"))
		clientRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Client")).Return(nil).Run(func(args mock.Arguments) {
			c := args.Get(1).(*models.Client)
			c.ID = "c-new"
		})

		c, err := svc.FindOrCreate(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "c-new", c.ID)
	})
}

func TestClientService_Update_AndDelete(t *testing.T) {
	clientRepo := new(mocks.MockClientRepository)
	appRepo := new(mocks.MockAppointmentRepository)

	svc := NewClientService(clientRepo, appRepo)

	t.Run("Update Notes", func(t *testing.T) {
		clientRepo.On("UpdateNotes", mock.Anything, "c1", "master-1", "Новая заметка").Return(nil)
		err := svc.UpdateNotes(context.Background(), "c1", "master-1", "Новая заметка")
		assert.NoError(t, err)
	})

	t.Run("Toggle Block", func(t *testing.T) {
		clientRepo.On("ToggleBlock", mock.Anything, "c1", "master-1", true).Return(nil)
		err := svc.ToggleBlock(context.Background(), "c1", "master-1", true)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		clientRepo.On("Delete", mock.Anything, "c1", "master-1").Return(nil)
		err := svc.Delete(context.Background(), "c1", "master-1")
		assert.NoError(t, err)
	})
}
