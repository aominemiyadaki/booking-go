package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-backend/internal/models"
	"booking-backend/internal/repository/mocks"
)

func TestLookupService(t *testing.T) {
	repo := new(mocks.MockLookupRepository)
	svc := NewLookupService(repo)

	ctx := context.Background()

	t.Run("Specializations", func(t *testing.T) {
		repo.On("GetSpecializations", ctx).Return([]models.Specialization{{ID: "sp1", Name: "Маникюр"}}, nil)
		res, err := svc.GetSpecializations(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		repo.On("CreateSpecialization", ctx, mock.AnythingOfType("*models.Specialization")).Return(nil)
		sp, err := svc.CreateSpecialization(ctx, "Педикюр")
		assert.NoError(t, err)
		assert.Equal(t, "Педикюр", sp.Name)

		repo.On("UpdateSpecialization", ctx, mock.AnythingOfType("*models.Specialization")).Return(nil)
		err = svc.UpdateSpecialization(ctx, "sp1", "Визаж")
		assert.NoError(t, err)

		repo.On("DeleteSpecialization", ctx, "sp1").Return(nil)
		err = svc.DeleteSpecialization(ctx, "sp1")
		assert.NoError(t, err)
	})

	t.Run("Currencies", func(t *testing.T) {
		repo.On("GetCurrencies", ctx).Return([]models.Currency{{ID: "c1", ISOCode: "RUB"}}, nil)
		res, err := svc.GetCurrencies(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		repo.On("CreateCurrency", ctx, mock.AnythingOfType("*models.Currency")).Return(nil)
		curr, err := svc.CreateCurrency(ctx, "Рубль", "RUB")
		assert.NoError(t, err)
		assert.Equal(t, "RUB", curr.ISOCode)

		repo.On("UpdateCurrency", ctx, mock.AnythingOfType("*models.Currency")).Return(nil)
		err = svc.UpdateCurrency(ctx, "c1", "Рубль", "RUB")
		assert.NoError(t, err)

		repo.On("DeleteCurrency", ctx, "c1").Return(nil)
		err = svc.DeleteCurrency(ctx, "c1")
		assert.NoError(t, err)
	})

	t.Run("AppointmentStatuses", func(t *testing.T) {
		repo.On("GetAppointmentStatuses", ctx).Return([]models.AppointmentStatus{{ID: "st1", Name: "подтвержден"}}, nil)
		res, err := svc.GetAppointmentStatuses(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		repo.On("CreateAppointmentStatus", ctx, mock.AnythingOfType("*models.AppointmentStatus")).Return(nil)
		st, err := svc.CreateAppointmentStatus(ctx, "новое")
		assert.NoError(t, err)
		assert.Equal(t, "новое", st.Name)

		repo.On("UpdateAppointmentStatus", ctx, mock.AnythingOfType("*models.AppointmentStatus")).Return(nil)
		err = svc.UpdateAppointmentStatus(ctx, "st1", "отменен")
		assert.NoError(t, err)

		repo.On("DeleteAppointmentStatus", ctx, "st1").Return(nil)
		err = svc.DeleteAppointmentStatus(ctx, "st1")
		assert.NoError(t, err)
	})

	t.Run("BreakTypes", func(t *testing.T) {
		repo.On("GetBreakTypes", ctx).Return([]models.BreakType{{ID: "bt1", Name: "Обед"}}, nil)
		res, err := svc.GetBreakTypes(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		repo.On("CreateBreakType", ctx, mock.AnythingOfType("*models.BreakType")).Return(nil)
		bt, err := svc.CreateBreakType(ctx, "Перерыв")
		assert.NoError(t, err)
		assert.Equal(t, "Перерыв", bt.Name)

		repo.On("UpdateBreakType", ctx, mock.AnythingOfType("*models.BreakType")).Return(nil)
		err = svc.UpdateBreakType(ctx, "bt1", "Перекур")
		assert.NoError(t, err)

		repo.On("DeleteBreakType", ctx, "bt1").Return(nil)
		err = svc.DeleteBreakType(ctx, "bt1")
		assert.NoError(t, err)
	})
}
