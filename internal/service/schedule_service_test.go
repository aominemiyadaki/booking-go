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

func TestScheduleService_GetAvailableSlots(t *testing.T) {
	schedRepo := new(mocks.MockScheduleRepository)
	svcRepo := new(mocks.MockServiceRepository)
	appRepo := new(mocks.MockAppointmentRepository)

	svc := NewScheduleService(schedRepo, svcRepo, appRepo)

	ctx := context.Background()
	futureDateStr := time.Now().AddDate(0, 0, 2).Format("2006-01-02")
	parsedDate, _ := time.ParseInLocation("2006-01-02", futureDateStr, time.Local)

	q := models.AvailableSlotsQuery{
		MasterID:  "master-1",
		Date:      futureDateStr,
		ServiceID: "service-1",
	}

	svcRepo.On("GetByID", ctx, "service-1").Return(&models.Service{
		ID:              "service-1",
		DurationMinutes: 60,
	}, nil)

	schedRepo.On("GetWorkingSlotsForDate", ctx, "master-1", futureDateStr).Return([]models.TimeSlot{
		{Start: "09:00", End: "18:00"},
	}, nil)

	appRepo.On("GetByDate", ctx, "master-1", mock.MatchedBy(func(t time.Time) bool {
		return t.Format("2006-01-02") == futureDateStr
	})).Return([]models.AppointmentResponse{
		{
			StartTime:  time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 12, 0, 0, 0, time.Local),
			EndTime:    time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 13, 0, 0, 0, time.Local),
			StatusName: "подтвержден",
		},
	}, nil)

	schedRepo.On("GetBreaks", ctx, "master-1").Return([]models.MasterBreak{}, nil)

	slots, err := svc.GetAvailableSlots(ctx, q)

	assert.NoError(t, err)
	assert.NotEmpty(t, slots)
	assert.NotContains(t, slots, "12:00") // Slot 12:00 is booked
	assert.Contains(t, slots, "09:00")
	assert.Contains(t, slots, "10:00")

	schedRepo.AssertExpectations(t)
	svcRepo.AssertExpectations(t)
	appRepo.AssertExpectations(t)
}

func TestScheduleService_SaveAndBreaks(t *testing.T) {
	schedRepo := new(mocks.MockScheduleRepository)
	svcRepo := new(mocks.MockServiceRepository)
	appRepo := new(mocks.MockAppointmentRepository)

	svc := NewScheduleService(schedRepo, svcRepo, appRepo)
	ctx := context.Background()

	t.Run("SaveTimeSlots", func(t *testing.T) {
		req := models.SaveTimeSlotsRequest{
			Dates: []string{"2026-08-15"},
			Slots: []models.TimeSlot{{Start: "10:00", End: "19:00"}},
		}
		schedRepo.On("SaveTimeSlots", ctx, "master-1", req.Dates, req.Slots).Return(nil)

		err := svc.SaveTimeSlots(ctx, "master-1", req)
		assert.NoError(t, err)
	})

	t.Run("GetWorkingDays", func(t *testing.T) {
		q := models.WorkingDaysQuery{MasterID: "master-1", Month: "2026-08"}
		schedRepo.On("GetWorkingDaysForMonth", ctx, "master-1", "2026-08").Return([]string{"2026-08-15"}, nil)

		days, err := svc.GetWorkingDays(ctx, q)
		assert.NoError(t, err)
		assert.Equal(t, []string{"2026-08-15"}, days)
	})

	t.Run("GetForMonth", func(t *testing.T) {
		schedRepo.On("GetForMonth", ctx, "master-1", "2026-08").Return(&models.MonthlySchedule{MasterID: "master-1"}, nil)

		ms, err := svc.GetForMonth(ctx, "master-1", "2026-08")
		assert.NoError(t, err)
		assert.NotNil(t, ms)
	})

	t.Run("DeleteWorkingDay", func(t *testing.T) {
		schedRepo.On("DeleteWorkingDay", ctx, "master-1", "2026-08-15").Return(nil)

		err := svc.DeleteWorkingDay(ctx, "master-1", "2026-08-15")
		assert.NoError(t, err)
	})

	t.Run("CreateBreak", func(t *testing.T) {
		req := models.CreateMasterBreakRequest{
			BreakTypeID: "bt-1",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(30 * time.Minute),
		}
		schedRepo.On("CreateBreak", ctx, mock.AnythingOfType("*models.MasterBreak")).Return(nil)

		res, err := svc.CreateBreak(ctx, "master-1", req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("GetBreaks", func(t *testing.T) {
		schedRepo.On("GetBreaks", ctx, "master-1").Return([]models.MasterBreak{{ID: "b1"}}, nil)

		breaks, err := svc.GetBreaks(ctx, "master-1")
		assert.NoError(t, err)
		assert.Len(t, breaks, 1)
	})

	t.Run("DeleteBreak", func(t *testing.T) {
		schedRepo.On("DeleteBreak", ctx, "b1", "master-1").Return(nil)

		err := svc.DeleteBreak(ctx, "b1", "master-1")
		assert.NoError(t, err)
	})
}
