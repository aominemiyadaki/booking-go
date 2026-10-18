package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"booking-backend/internal/models"
)

// MockMasterRepository — мок для MasterRepositoryInterface
type MockMasterRepository struct {
	mock.Mock
}

func (m *MockMasterRepository) Create(ctx context.Context, master *models.Master) error {
	args := m.Called(ctx, master)
	return args.Error(0)
}

func (m *MockMasterRepository) GetByEmail(ctx context.Context, email string) (*models.Master, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Master), args.Error(1)
}

func (m *MockMasterRepository) GetByID(ctx context.Context, id string) (*models.Master, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Master), args.Error(1)
}

func (m *MockMasterRepository) GetAll(ctx context.Context) ([]models.MasterPublic, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.MasterPublic), args.Error(1)
}

func (m *MockMasterRepository) Update(ctx context.Context, master *models.Master) error {
	args := m.Called(ctx, master)
	return args.Error(0)
}

func (m *MockMasterRepository) UpdateSessionToken(ctx context.Context, masterID, token string) error {
	args := m.Called(ctx, masterID, token)
	return args.Error(0)
}

// MockServiceRepository — мок для ServiceRepositoryInterface
type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) GetByMasterID(ctx context.Context, masterID string) ([]models.Service, error) {
	args := m.Called(ctx, masterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Service), args.Error(1)
}

func (m *MockServiceRepository) GetByID(ctx context.Context, id string) (*models.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Service), args.Error(1)
}

func (m *MockServiceRepository) Create(ctx context.Context, s *models.Service) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockServiceRepository) Update(ctx context.Context, s *models.Service) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockServiceRepository) Delete(ctx context.Context, id, masterID string) error {
	args := m.Called(ctx, id, masterID)
	return args.Error(0)
}

// MockClientRepository — мок для ClientRepositoryInterface
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) GetByMasterID(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.Client, int64, error) {
	args := m.Called(ctx, masterID, page)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Client), args.Get(1).(int64), args.Error(2)
}

func (m *MockClientRepository) GetByID(ctx context.Context, id string) (*models.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientRepository) GetByIDAndMaster(ctx context.Context, id, masterID string) (*models.Client, error) {
	args := m.Called(ctx, id, masterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientRepository) GetByPhone(ctx context.Context, phone string) (*models.Client, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientRepository) Create(ctx context.Context, c *models.Client) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockClientRepository) Update(ctx context.Context, c *models.Client) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockClientRepository) Delete(ctx context.Context, id, masterID string) error {
	args := m.Called(ctx, id, masterID)
	return args.Error(0)
}

func (m *MockClientRepository) UpdateNotes(ctx context.Context, id, masterID, notes string) error {
	args := m.Called(ctx, id, masterID, notes)
	return args.Error(0)
}

func (m *MockClientRepository) ToggleBlock(ctx context.Context, id, masterID string, block bool) error {
	args := m.Called(ctx, id, masterID, block)
	return args.Error(0)
}

func (m *MockClientRepository) IncrementVisits(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockAppointmentRepository — мок для AppointmentRepositoryInterface
type MockAppointmentRepository struct {
	mock.Mock
}

func (m *MockAppointmentRepository) Create(ctx context.Context, a *models.Appointment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAppointmentRepository) GetByID(ctx context.Context, id string) (*models.Appointment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) GetByMasterID(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error) {
	args := m.Called(ctx, masterID, page)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.AppointmentResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockAppointmentRepository) GetByClientID(ctx context.Context, clientID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error) {
	args := m.Called(ctx, clientID, page)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.AppointmentResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockAppointmentRepository) GetByDate(ctx context.Context, masterID string, date time.Time) ([]models.AppointmentResponse, error) {
	args := m.Called(ctx, masterID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AppointmentResponse), args.Error(1)
}

func (m *MockAppointmentRepository) GetAppointmentResponseByID(ctx context.Context, id string) (*models.AppointmentResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AppointmentResponse), args.Error(1)
}

func (m *MockAppointmentRepository) Update(ctx context.Context, a *models.Appointment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAppointmentRepository) UpdateStatus(ctx context.Context, id, masterID, statusID string) error {
	args := m.Called(ctx, id, masterID, statusID)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Delete(ctx context.Context, id, masterID string) error {
	args := m.Called(ctx, id, masterID)
	return args.Error(0)
}

// MockBlacklistRepository — мок для BlacklistRepositoryInterface
type MockBlacklistRepository struct {
	mock.Mock
}

func (m *MockBlacklistRepository) GetByMasterID(ctx context.Context, masterID string) ([]models.BlacklistedClient, error) {
	args := m.Called(ctx, masterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BlacklistedClient), args.Error(1)
}

func (m *MockBlacklistRepository) Add(ctx context.Context, b *models.BlacklistedClient) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBlacklistRepository) Remove(ctx context.Context, id, masterID string) error {
	args := m.Called(ctx, id, masterID)
	return args.Error(0)
}

func (m *MockBlacklistRepository) IsBlocked(ctx context.Context, masterID, clientID string) (bool, error) {
	args := m.Called(ctx, masterID, clientID)
	return args.Bool(0), args.Error(1)
}

// MockScheduleRepository — мок для ScheduleRepositoryInterface
type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) GetWorkingSlotsForDate(ctx context.Context, masterID, date string) ([]models.TimeSlot, error) {
	args := m.Called(ctx, masterID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TimeSlot), args.Error(1)
}

func (m *MockScheduleRepository) SaveTimeSlots(ctx context.Context, masterID string, dates []string, slots []models.TimeSlot) error {
	args := m.Called(ctx, masterID, dates, slots)
	return args.Error(0)
}

func (m *MockScheduleRepository) DeleteWorkingDay(ctx context.Context, masterID, date string) error {
	args := m.Called(ctx, masterID, date)
	return args.Error(0)
}

func (m *MockScheduleRepository) GetWorkingDaysForMonth(ctx context.Context, masterID, yearMonth string) ([]string, error) {
	args := m.Called(ctx, masterID, yearMonth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockScheduleRepository) GetForMonth(ctx context.Context, masterID, yearMonth string) (*models.MonthlySchedule, error) {
	args := m.Called(ctx, masterID, yearMonth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MonthlySchedule), args.Error(1)
}

func (m *MockScheduleRepository) CreateBreak(ctx context.Context, b *models.MasterBreak) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockScheduleRepository) GetBreaks(ctx context.Context, masterID string) ([]models.MasterBreak, error) {
	args := m.Called(ctx, masterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.MasterBreak), args.Error(1)
}

func (m *MockScheduleRepository) DeleteBreak(ctx context.Context, id, masterID string) error {
	args := m.Called(ctx, id, masterID)
	return args.Error(0)
}

// MockLookupRepository — мок для LookupRepositoryInterface
type MockLookupRepository struct {
	mock.Mock
}

func (m *MockLookupRepository) GetSpecializations(ctx context.Context) ([]models.Specialization, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Specialization), args.Error(1)
}

func (m *MockLookupRepository) CreateSpecialization(ctx context.Context, s *models.Specialization) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockLookupRepository) UpdateSpecialization(ctx context.Context, s *models.Specialization) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockLookupRepository) DeleteSpecialization(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLookupRepository) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Currency), args.Error(1)
}

func (m *MockLookupRepository) CreateCurrency(ctx context.Context, c *models.Currency) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockLookupRepository) UpdateCurrency(ctx context.Context, c *models.Currency) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockLookupRepository) DeleteCurrency(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLookupRepository) GetAppointmentStatuses(ctx context.Context) ([]models.AppointmentStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AppointmentStatus), args.Error(1)
}

func (m *MockLookupRepository) CreateAppointmentStatus(ctx context.Context, s *models.AppointmentStatus) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockLookupRepository) UpdateAppointmentStatus(ctx context.Context, s *models.AppointmentStatus) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockLookupRepository) DeleteAppointmentStatus(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLookupRepository) GetBreakTypes(ctx context.Context) ([]models.BreakType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BreakType), args.Error(1)
}

func (m *MockLookupRepository) CreateBreakType(ctx context.Context, b *models.BreakType) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockLookupRepository) UpdateBreakType(ctx context.Context, b *models.BreakType) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockLookupRepository) DeleteBreakType(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
