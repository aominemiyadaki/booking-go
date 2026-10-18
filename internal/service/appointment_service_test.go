package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-backend/internal/models"
	"booking-backend/internal/repository/mocks"
)

func TestAppointmentService_Create_Success(t *testing.T) {
	mockAppRepo := new(mocks.MockAppointmentRepository)
	mockClientRepo := new(mocks.MockClientRepository)
	mockBlacklistRepo := new(mocks.MockBlacklistRepository)
	mockServiceRepo := new(mocks.MockServiceRepository)
	mockLookupRepo := new(mocks.MockLookupRepository)

	svc := NewAppointmentService(mockAppRepo, mockClientRepo, mockBlacklistRepo, mockServiceRepo, mockLookupRepo)

	req := models.BookAppointmentRequest{
		MasterID:  "master-1",
		ClientID:  "client-1",
		ServiceID: "service-1",
		Date:      "2026-07-12",
		SlotTime:  "10:00",
		StatusID:  "status-1",
	}

	mockBlacklistRepo.On("IsBlocked", mock.Anything, "master-1", "client-1").Return(false, nil)
	
	mockServiceRepo.On("GetByID", mock.Anything, "service-1").Return(&models.Service{
		ID:              "service-1",
		DurationMinutes: 60,
	}, nil)

	mockLookupRepo.On("GetAppointmentStatuses", mock.Anything).Return([]models.AppointmentStatus{
		{ID: "status-1", Name: "в ожидании"},
	}, nil)

	var createdApp *models.Appointment
	mockAppRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Appointment")).Return(nil).Run(func(args mock.Arguments) {
		createdApp = args.Get(1).(*models.Appointment)
	})

	mockAppRepo.On("GetAppointmentResponseByID", mock.Anything, mock.Anything).Return(&models.AppointmentResponse{
		MasterID:  "master-1",
		ClientID:  "client-1",
		StartTime: time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 7, 12, 11, 0, 0, 0, time.UTC),
	}, nil)

	ctx := context.WithValue(context.Background(), "request_id", "test-trace-id-123")

	appointment, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, appointment)
	assert.NotNil(t, createdApp)
	assert.Equal(t, 10, createdApp.StartTime.Hour())

	mockBlacklistRepo.AssertExpectations(t)
	mockServiceRepo.AssertExpectations(t)
	mockAppRepo.AssertExpectations(t)
}

func TestAppointmentService_Create_Blacklisted(t *testing.T) {
	mockAppRepo := new(mocks.MockAppointmentRepository)
	mockClientRepo := new(mocks.MockClientRepository)
	mockBlacklistRepo := new(mocks.MockBlacklistRepository)
	mockServiceRepo := new(mocks.MockServiceRepository)
	mockLookupRepo := new(mocks.MockLookupRepository)

	svc := NewAppointmentService(mockAppRepo, mockClientRepo, mockBlacklistRepo, mockServiceRepo, mockLookupRepo)

	req := models.BookAppointmentRequest{
		MasterID:  "master-1",
		ClientID:  "client-1",
		ServiceID: "service-1",
		Date:      "2026-07-12",
		SlotTime:  "10:00",
		StatusID:  "status-1",
	}

	mockBlacklistRepo.On("IsBlocked", mock.Anything, "master-1", "client-1").Return(true, nil)

	ctx := context.WithValue(context.Background(), "request_id", "test-trace-id-456")

	_, err := svc.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "онлайн-запись для вас ограничена", err.Error())
}

func TestAppointmentService_GetByID(t *testing.T) {
	mockAppRepo := new(mocks.MockAppointmentRepository)
	svc := NewAppointmentService(mockAppRepo, nil, nil, nil, nil)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockAppRepo.On("GetByID", ctx, "app-1").Return(&models.Appointment{ID: "app-1", MasterID: "master-1"}, nil)
		res, err := svc.GetByID(ctx, "app-1", "master-1")
		assert.NoError(t, err)
		assert.Equal(t, "app-1", res.ID)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockAppRepo.On("GetByID", ctx, "app-2").Return(&models.Appointment{ID: "app-2", MasterID: "master-other"}, nil)
		_, err := svc.GetByID(ctx, "app-2", "master-1")
		assert.Error(t, err)
		assert.Equal(t, "доступ запрещён: запись принадлежит другому мастеру", err.Error())
	})
}

func TestAppointmentService_Update(t *testing.T) {
	mockAppRepo := new(mocks.MockAppointmentRepository)
	mockSvcRepo := new(mocks.MockServiceRepository)
	svc := NewAppointmentService(mockAppRepo, nil, nil, mockSvcRepo, nil)
	ctx := context.Background()

	mockAppRepo.On("GetByID", ctx, "app-1").Return(&models.Appointment{ID: "app-1", MasterID: "master-1"}, nil)
	mockSvcRepo.On("GetByID", ctx, "s1").Return(&models.Service{ID: "s1", DurationMinutes: 30}, nil)
	mockAppRepo.On("Update", ctx, mock.AnythingOfType("*models.Appointment")).Return(nil)

	err := svc.Update(ctx, "app-1", "master-1", models.BookAppointmentRequest{
		ServiceID: "s1",
		Date:      "2026-08-20",
		SlotTime:  "14:00",
	})

	assert.NoError(t, err)
}

func TestAppointmentService_UpdateStatus(t *testing.T) {
	mockAppRepo := new(mocks.MockAppointmentRepository)
	mockClientRepo := new(mocks.MockClientRepository)
	mockLookupRepo := new(mocks.MockLookupRepository)

	svc := NewAppointmentService(mockAppRepo, mockClientRepo, nil, nil, mockLookupRepo)
	ctx := context.WithValue(context.Background(), "request_id", "test-trace-id")

	mockLookupRepo.On("GetAppointmentStatuses", mock.Anything).Return([]models.AppointmentStatus{
		{ID: "status-1", Name: "в ожидании"},
		{ID: "status-2", Name: "подтвержден"},
	}, nil)

	mockAppRepo.On("GetByID", mock.Anything, "app-1").Return(&models.Appointment{
		ID:       "app-1",
		MasterID: "master-1",
		ClientID: "client-1",
		StatusID: "status-1",
	}, nil)

	mockAppRepo.On("UpdateStatus", mock.Anything, "app-1", "master-1", "status-2").Return(nil)
	mockClientRepo.On("IncrementVisits", mock.Anything, "client-1").Return(nil)

	err := svc.UpdateStatus(ctx, "app-1", "master-1", "подтвержден")

	assert.NoError(t, err)
}

func TestAppointmentService_GetGuestAppointments(t *testing.T) {
	mockAppRepo := new(mocks.MockAppointmentRepository)
	mockClientRepo := new(mocks.MockClientRepository)
	svc := NewAppointmentService(mockAppRepo, mockClientRepo, nil, nil, nil)
	ctx := context.Background()

	t.Run("By Phone Found", func(t *testing.T) {
		mockClientRepo.On("GetByPhone", ctx, "79991112233").Return(&models.Client{ID: "c1"}, nil)
		mockAppRepo.On("GetByClientID", ctx, "c1", models.PaginationRequest{Page: 1, PageSize: 100}).Return([]models.AppointmentResponse{{ID: "a1"}}, int64(1), nil)

		res, err := svc.GetGuestAppointments(ctx, "", "+7 (999) 111-22-33")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("By Guest Token", func(t *testing.T) {
		mockAppRepo.On("GetByID", ctx, "token-app-id").Return(&models.Appointment{ID: "token-app-id", ClientID: "c2"}, nil)
		mockAppRepo.On("GetByClientID", ctx, "c2", models.PaginationRequest{Page: 1, PageSize: 100}).Return([]models.AppointmentResponse{{ID: "token-app-id"}}, int64(1), nil)

		res, err := svc.GetGuestAppointments(ctx, "token-app-id", "")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestAppointmentService_DeleteAndMisc(t *testing.T) {
	mockAppRepo := new(mocks.MockAppointmentRepository)
	svc := NewAppointmentService(mockAppRepo, nil, nil, nil, nil)
	ctx := context.Background()

	mockAppRepo.On("Delete", ctx, "a1", "m1").Return(nil)
	err := svc.Delete(ctx, "a1", "m1")
	assert.NoError(t, err)

	mockAppRepo.On("GetByDate", ctx, "m1", mock.Anything).Return([]models.AppointmentResponse{}, nil)
	resDate, err := svc.GetByDate(ctx, "m1", time.Now())
	assert.NoError(t, err)
	assert.Empty(t, resDate)

	_, err = svc.CreateBlock(ctx, "m1", nil)
	assert.Error(t, err)
}
