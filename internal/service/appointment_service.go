package service

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"booking-backend/internal/logger"
	"booking-backend/internal/models"
	"booking-backend/internal/repository"
)

type AppointmentService struct {
	appointmentRepo repository.AppointmentRepositoryInterface
	clientRepo      repository.ClientRepositoryInterface
	blacklistRepo   repository.BlacklistRepositoryInterface
	serviceRepo     repository.ServiceRepositoryInterface
	lookupRepo      repository.LookupRepositoryInterface
}

func NewAppointmentService(
	appointmentRepo repository.AppointmentRepositoryInterface,
	clientRepo repository.ClientRepositoryInterface,
	blacklistRepo repository.BlacklistRepositoryInterface,
	serviceRepo repository.ServiceRepositoryInterface,
	lookupRepo repository.LookupRepositoryInterface,
) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		clientRepo:      clientRepo,
		blacklistRepo:   blacklistRepo,
		serviceRepo:     serviceRepo,
		lookupRepo:      lookupRepo,
	}
}

// Create создаёт новую запись клиента.
func (s *AppointmentService) Create(ctx context.Context, req models.BookAppointmentRequest) (*models.AppointmentResponse, error) {
	defer logger.Trace(ctx, "AppointmentService.Create", slog.String("master_id", req.MasterID), slog.String("client_id", req.ClientID))()

	reqID, _ := ctx.Value("request_id").(string)

	// 0. Если передан номер телефона, ищем клиента в базе по телефону. Если есть — берем его ID, если нет — создаем клиента.
	if req.ClientPhone != "" {
		normalizedPhone := normalizePhone(req.ClientPhone)
		client, err := s.clientRepo.GetByPhone(ctx, normalizedPhone)
		if err == nil && client != nil {
			req.ClientID = client.ID
		} else if req.ClientID == "" {
			if req.ClientName == "" {
				return nil, errors.New("для создания нового клиента укажите имя (client_name)")
			}
			newClient := &models.Client{
				MasterID:    req.MasterID,
				Name:        req.ClientName,
				PhoneNumber: normalizedPhone,
			}
			if err := s.clientRepo.Create(ctx, newClient); err != nil {
				return nil, errors.New("не удалось создать профиль клиента")
			}
			req.ClientID = newClient.ID
		}
	} else if req.ClientID == "" {
		return nil, errors.New("для записи нужен client_id или (client_name и client_phone)")
	}

	// 1. Проверяем blacklist
	blocked, err := s.blacklistRepo.IsBlocked(ctx, req.MasterID, req.ClientID)
	if err != nil {
		slog.Error("Ошибка проверки блэклиста", "error", err, "request_id", reqID)
		return nil, errors.New("ошибка сервера")
	}
	if blocked {
		slog.Warn("Заблокированный клиент пытается записаться", "client_id", req.ClientID, "request_id", reqID)
		return nil, errors.New("онлайн-запись для вас ограничена")
	}

	// 2. Получаем длительность услуги
	svc, err := s.serviceRepo.GetByID(ctx, req.ServiceID)
	if err != nil {
		return nil, errors.New("услуга не найдена")
	}

	// 3. Высчитываем время
	parsedDate, err := time.ParseInLocation("2006-01-02", req.Date, time.Local)
	if err != nil {
		return nil, errors.New("неверный формат даты, ожидается YYYY-MM-DD")
	}
	parsedTime, err := time.ParseInLocation("15:04", req.SlotTime, time.Local)
	if err != nil {
		return nil, errors.New("неверный формат времени, ожидается HH:MM")
	}

	startTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, parsedDate.Location())
	endTime := startTime.Add(time.Duration(svc.DurationMinutes) * time.Minute)

	// 4. Проверяем StatusID, если пустой - ищем статус "в ожидании"
	if req.StatusID == "" {
		statuses, err := s.lookupRepo.GetAppointmentStatuses(ctx)
		if err == nil {
			for _, st := range statuses {
				if st.Name == "в ожидании" {
					req.StatusID = st.ID
					break
				}
			}
		}
		// Если не нашли статус в БД, то база может выдать ошибку (так как поле NOT NULL)
	}

	// 5. Создаём запись
	appointment := &models.Appointment{
		MasterID:  req.MasterID,
		ClientID:  req.ClientID,
		ServiceID: req.ServiceID,
		StartTime: startTime,
		EndTime:   endTime,
		StatusID:  req.StatusID,
		Comment:   req.Comment,
	}

	if err := s.appointmentRepo.Create(ctx, appointment); err != nil {
		slog.Error("Ошибка создания записи в БД", "error", err, "request_id", reqID)
		return nil, errors.New("не удалось создать запись")
	}

	slog.Info("Запись успешно создана", "appointment_id", appointment.ID, "request_id", reqID)

	return s.appointmentRepo.GetAppointmentResponseByID(ctx, appointment.ID)
}

// GetByID находит запись по ID и проверяет права мастера.
func (s *AppointmentService) GetByID(ctx context.Context, id, masterID string) (*models.Appointment, error) {
	defer logger.Trace(ctx, "AppointmentService.GetByID", slog.String("appointment_id", id))()

	app, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app.MasterID != masterID {
		return nil, errors.New("доступ запрещён: запись принадлежит другому мастеру")
	}
	return app, nil
}

// Update обновляет данные записи (перенос времени, изменение комментариев и т.д.).
func (s *AppointmentService) Update(ctx context.Context, id, masterID string, req models.BookAppointmentRequest) error {
	defer logger.Trace(ctx, "AppointmentService.Update", slog.String("appointment_id", id))()

	reqID, _ := ctx.Value("request_id").(string)

	app, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if app.MasterID != masterID {
		slog.Warn("Попытка несанкционированного изменения записи другого мастера", "id", id, "master_id", masterID, "request_id", reqID)
		return errors.New("доступ запрещён: запись принадлежит другому мастеру")
	}

	svc, err := s.serviceRepo.GetByID(ctx, req.ServiceID)
	if err != nil {
		return errors.New("услуга не найдена")
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", req.Date, time.Local)
	if err != nil {
		return errors.New("неверный формат даты, ожидается YYYY-MM-DD")
	}
	parsedTime, err := time.ParseInLocation("15:04", req.SlotTime, time.Local)
	if err != nil {
		return errors.New("неверный формат времени, ожидается HH:MM")
	}

	startTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, parsedDate.Location()).UTC()
	endTime := startTime.Add(time.Duration(svc.DurationMinutes) * time.Minute)

	// Обновляем поля
	app.ServiceID = req.ServiceID
	app.StartTime = startTime
	app.EndTime = endTime
	app.StatusID = req.StatusID
	app.Comment = req.Comment

	if err := s.appointmentRepo.Update(ctx, app); err != nil {
		slog.Error("Не удалось обновить запись в БД", "error", err, "request_id", reqID)
		return err
	}

	slog.Info("Запись успешно обновлена", "appointment_id", id, "request_id", reqID)
	return nil
}

// GetByMasterID возвращает записи мастера.
func (s *AppointmentService) GetByMasterID(ctx context.Context, masterID string, page models.PaginationRequest) ([]models.AppointmentResponse, int64, error) {
	return s.appointmentRepo.GetByMasterID(ctx, masterID, page)
}

// UpdateStatus изменяет статус записи.
func (s *AppointmentService) UpdateStatus(ctx context.Context, id, masterID string, statusStr string) error {
	reqID, _ := ctx.Value("request_id").(string)

	statusID := statusStr
	switch statusStr {
	case "pending":
		statusStr = "в ожидании"
	case "confirmed":
		statusStr = "подтвержден"
	case "cancelled":
		statusStr = "отменен"
	case "completed":
		statusStr = "завершен"
	}
	statuses, err := s.lookupRepo.GetAppointmentStatuses(ctx)
	if err == nil {
		for _, st := range statuses {
			if st.Name == statusStr || st.ID == statusStr {
				statusID = st.ID
				break
			}
		}
	}

	app, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.appointmentRepo.UpdateStatus(ctx, id, masterID, statusID)
	if err != nil {
		slog.Error("Ошибка обновления статуса", "error", err, "request_id", reqID)
		return err
	}

	if app.StatusID != statusID && statusStr == "подтвержден" {
		if incErr := s.clientRepo.IncrementVisits(ctx, app.ClientID); incErr != nil {
			slog.Warn("Не удалось обновить счетчик посещений", "client_id", app.ClientID, "error", incErr, "request_id", reqID)
		}
	}

	return nil
}

// Delete удаляет запись.
func (s *AppointmentService) Delete(ctx context.Context, id, masterID string) error {
	return s.appointmentRepo.Delete(ctx, id, masterID)
}

// GetByDate возвращает записи мастера на конкретную дату.
func (s *AppointmentService) GetByDate(ctx context.Context, masterID string, date time.Time) ([]models.AppointmentResponse, error) {
	return s.appointmentRepo.GetByDate(ctx, masterID, date)
}

// GetGuestAppointments возвращает записи гостя по его телефону или токену (ID последней записи).
func (s *AppointmentService) GetGuestAppointments(ctx context.Context, guestToken, phone string) ([]models.AppointmentResponse, error) {
	reqID, _ := ctx.Value("request_id").(string)

	var clientID string

	if phone != "" {
		normalized := normalizePhone(phone)
		client, err := s.clientRepo.GetByPhone(ctx, normalized)
		if err != nil {
			slog.Warn("Гость с таким телефоном не найден в базе клиентов", "phone", phone, "request_id", reqID)
			return []models.AppointmentResponse{}, nil
		}
		clientID = client.ID
	} else if guestToken != "" {
		app, err := s.appointmentRepo.GetByID(ctx, guestToken)
		if err != nil {
			slog.Warn("Запись по guestToken (appointment ID) не найдена", "guest_token", guestToken, "request_id", reqID)
			return []models.AppointmentResponse{}, nil
		}
		clientID = app.ClientID
	} else {
		return nil, errors.New("параметр guest_token или phone обязателен")
	}

	page := models.PaginationRequest{Page: 1, PageSize: 100}
	appointments, _, err := s.appointmentRepo.GetByClientID(ctx, clientID, page)
	return appointments, err
}

// CreateBlock — заглушка для обратной совместимости создания перерывов через старый роут.
func (s *AppointmentService) CreateBlock(ctx context.Context, masterID string, req interface{}) (*models.Appointment, error) {
	return nil, errors.New("создание блоков через данный роут отключено, используйте расписание")
}

var nonDigitRegex = regexp.MustCompile(`\D`)

func normalizePhone(phone string) string {
	return nonDigitRegex.ReplaceAllString(phone, "")
}
