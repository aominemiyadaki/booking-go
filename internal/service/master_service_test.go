package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"booking-backend/internal/models"
	"booking-backend/internal/repository/mocks"
)

func TestMasterService_Register(t *testing.T) {
	tests := []struct {
		name        string
		req         models.RegisterRequest
		setupMock   func(repo *mocks.MockMasterRepository)
		expectedErr string
	}{
		{
			name: "Success",
			req: models.RegisterRequest{
				Email:       "test@example.com",
				Password:    "password123",
				FullName:    "Иван Иванов",
				PhoneNumber: "+79001112233",
			},
			setupMock: func(repo *mocks.MockMasterRepository) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("not found"))
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Master")).Return(nil).Run(func(args mock.Arguments) {
					m := args.Get(1).(*models.Master)
					m.ID = "master-uuid-1"
				})
				repo.On("GetByID", mock.Anything, "master-uuid-1").Return(&models.Master{
					ID:                 "master-uuid-1",
					Email:              "test@example.com",
					FullName:           "Иван Иванов",
					PhoneNumber:        "+79001112233",
					SpecializationName: "Парикмахер",
				}, nil)
			},
			expectedErr: "",
		},
		{
			name: "User Already Exists",
			req: models.RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func(repo *mocks.MockMasterRepository) {
				repo.On("GetByEmail", mock.Anything, "existing@example.com").Return(&models.Master{ID: "master-existing"}, nil)
			},
			expectedErr: "пользователь с таким email уже существует",
		},
		{
			name: "Password Too Short",
			req: models.RegisterRequest{
				Email:    "new@example.com",
				Password: "short",
			},
			setupMock: func(repo *mocks.MockMasterRepository) {
				repo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, errors.New("not found"))
			},
			expectedErr: "пароль должен содержать минимум 8 символов",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockMasterRepository)
			tt.setupMock(repo)

			svc := NewMasterService(repo, "test-secret-key-32-bytes-long!!")
			res, err := svc.Register(context.Background(), tt.req)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.Token)
				assert.Equal(t, "Иван Иванов", res.Master.FullName)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestMasterService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("validpassword"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		req         models.LoginRequest
		setupMock   func(repo *mocks.MockMasterRepository)
		expectedErr string
	}{
		{
			name: "Success",
			req: models.LoginRequest{
				Email:    "master@example.com",
				Password: "validpassword",
			},
			setupMock: func(repo *mocks.MockMasterRepository) {
				repo.On("GetByEmail", mock.Anything, "master@example.com").Return(&models.Master{
					ID:           "master-1",
					Email:        "master@example.com",
					PasswordHash: string(hashedPassword),
					FullName:     "Мастер",
				}, nil)
			},
			expectedErr: "",
		},
		{
			name: "Email Not Found",
			req: models.LoginRequest{
				Email:    "unknown@example.com",
				Password: "password",
			},
			setupMock: func(repo *mocks.MockMasterRepository) {
				repo.On("GetByEmail", mock.Anything, "unknown@example.com").Return(nil, errors.New("not found"))
			},
			expectedErr: "неверный email или пароль",
		},
		{
			name: "Wrong Password",
			req: models.LoginRequest{
				Email:    "master@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(repo *mocks.MockMasterRepository) {
				repo.On("GetByEmail", mock.Anything, "master@example.com").Return(&models.Master{
					ID:           "master-1",
					Email:        "master@example.com",
					PasswordHash: string(hashedPassword),
				}, nil)
			},
			expectedErr: "неверный email или пароль",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockMasterRepository)
			tt.setupMock(repo)

			svc := NewMasterService(repo, "test-secret-key-32-bytes-long!!")
			res, err := svc.Login(context.Background(), tt.req)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.Token)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestMasterService_GetProfile_And_Masters(t *testing.T) {
	repo := new(mocks.MockMasterRepository)
	ctx := context.Background()

	t.Run("GetProfile Success", func(t *testing.T) {
		repo.On("GetByID", ctx, "master-1").Return(&models.Master{
			ID:          "master-1",
			FullName:    "Анна Смирнова",
			PhoneNumber: "+79998887766",
		}, nil)

		svc := NewMasterService(repo, "secret")
		profile, err := svc.GetProfile(ctx, "master-1")

		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, "Анна Смирнова", profile.FullName)
	})

	t.Run("GetMasters", func(t *testing.T) {
		repo.On("GetAll", ctx).Return([]models.MasterPublic{{ID: "m1", FullName: "Мастер 1"}}, nil)
		svc := NewMasterService(repo, "secret")
		masters, err := svc.GetMasters(ctx)

		assert.NoError(t, err)
		assert.Len(t, masters, 1)
	})

	t.Run("UpdateSessionToken", func(t *testing.T) {
		repo.On("UpdateSessionToken", ctx, "m1", "token123").Return(nil)
		svc := NewMasterService(repo, "secret")
		err := svc.UpdateSessionToken(ctx, "m1", "token123")

		assert.NoError(t, err)
	})
}

func TestMasterService_UpdateProfile(t *testing.T) {
	repo := new(mocks.MockMasterRepository)
	repo.On("GetByID", mock.Anything, "master-1").Return(&models.Master{
		ID:       "master-1",
		FullName: "Старое Имя",
	}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*models.Master")).Return(nil)

	svc := NewMasterService(repo, "secret")
	err := svc.UpdateProfile(context.Background(), "master-1", models.UpdateMasterProfileRequest{
		ID:       "master-1",
		FullName: "Новое Имя",
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
