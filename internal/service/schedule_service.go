package service

import (
	"context"
	"time"

	"booking-backend/internal/models"
	"booking-backend/internal/repository"
)

type ScheduleService struct {
	repo    repository.ScheduleRepositoryInterface
	svcRepo repository.ServiceRepositoryInterface
	appRepo repository.AppointmentRepositoryInterface
}

func NewScheduleService(
	repo repository.ScheduleRepositoryInterface,
	svcRepo repository.ServiceRepositoryInterface,
	appRepo repository.AppointmentRepositoryInterface,
) *ScheduleService {
	return &ScheduleService{
		repo:    repo,
		svcRepo: svcRepo,
		appRepo: appRepo,
	}
}

// Вспомогательная структура для расчета свободных окон
type timePeriod struct {
	Start time.Time
	End   time.Time
}

func subtractPeriod(busy timePeriod, periods []timePeriod) []timePeriod {
	// 1. Защита от дурака: если занятый период невалиден, ничего не вычитаем
	if !busy.Start.Before(busy.End) {
		return periods
	}

	// 2. Предварительная аллокация памяти для оптимизации (capacity = len(periods))
	result := make([]timePeriod, 0, len(periods))

	for _, p := range periods {
		// 3. Классическая формула проверки пересечения отрезков
		if busy.Start.Before(p.End) && busy.End.After(p.Start) {

			// Есть пересечение. Откусываем нужные части.
			if busy.Start.After(p.Start) {
				result = append(result, timePeriod{Start: p.Start, End: busy.Start})
			}
			if busy.End.Before(p.End) {
				result = append(result, timePeriod{Start: busy.End, End: p.End})
			}
		} else {
			// Пересечений нет, сохраняем исходный период
			result = append(result, p)
		}
	}

	return result
}

func (s *ScheduleService) GetAvailableSlots(ctx context.Context, q models.AvailableSlotsQuery) ([]string, error) {
	// 1. По умолчанию шаг и длительность 30 минут
	duration := 30 * time.Minute
	step := 30 * time.Minute

	// Если передана услуга, берем ее длительность
	if q.ServiceID != "" {
		svc, err := s.svcRepo.GetByID(ctx, q.ServiceID)
		if err == nil {
			duration = time.Duration(svc.DurationMinutes) * time.Minute
		}
	}

	// 2. Получаем рабочие слоты мастера на этот день
	workingSlots, err := s.repo.GetWorkingSlotsForDate(ctx, q.MasterID, q.Date)
	if err != nil {
		return nil, err
	}
	if len(workingSlots) == 0 {
		return []string{}, nil // Выходной
	}

	// Парсим дату запроса в локальной зоне
	parsedDate, err := time.ParseInLocation("2006-01-02", q.Date, time.Local)
	if err != nil {
		return nil, err
	}

	// 3. Формируем изначальные "свободные окна" из рабочих слотов
	var freePeriods []timePeriod
	for _, ws := range workingSlots {
		// Время в БД может быть "15:04:05" или "15:04"
		startT, err := time.Parse("15:04:05", ws.Start)
		if err != nil {
			startT, _ = time.Parse("15:04", ws.Start)
		}
		endT, err := time.Parse("15:04:05", ws.End)
		if err != nil {
			endT, _ = time.Parse("15:04", ws.End)
		}

		// Комбинируем дату из запроса и время из слота
		startDateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), startT.Hour(), startT.Minute(), startT.Second(), 0, parsedDate.Location())
		endDateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), endT.Hour(), endT.Minute(), endT.Second(), 0, parsedDate.Location())

		freePeriods = append(freePeriods, timePeriod{Start: startDateTime, End: endDateTime})
	}

	// 4. Получаем все существующие записи на этот день и вычитаем их
	appointments, err := s.appRepo.GetByDate(ctx, q.MasterID, parsedDate)
	if err != nil {
		return nil, err
	}
	for _, app := range appointments {
		if app.StatusName == "отменен" {
			continue
		}
		busy := timePeriod{Start: app.StartTime, End: app.EndTime}
		freePeriods = subtractPeriod(busy, freePeriods)
	}

	// 5. Получаем перерывы мастера и вычитаем их (фильтруем по текущей дате)
	breaks, err := s.repo.GetBreaks(ctx, q.MasterID)
	if err != nil {
		return nil, err
	}
	for _, b := range breaks {
		if b.StartTime.Year() == parsedDate.Year() && b.StartTime.Month() == parsedDate.Month() && b.StartTime.Day() == parsedDate.Day() {
			busy := timePeriod{Start: b.StartTime, End: b.EndTime}
			freePeriods = subtractPeriod(busy, freePeriods)
		}
	}

	// 6. Нарезка слотов с шагом 30 минут
	var result []string

	now := time.Now()

	for _, p := range freePeriods {
		current := p.Start
		// Идем по окну с шагом в 30 минут, пока длительность помещается до конца окна
		for current.Add(duration).Before(p.End) || current.Add(duration).Equal(p.End) {
			// Добавляем слот только если он в будущем (чтобы скрыть прошедшие слоты для текущего дня)
			if current.After(now) {
				result = append(result, current.Format("15:04"))
			}
			current = current.Add(step)
		}
	}

	return result, nil
}

func (s *ScheduleService) GetWorkingDays(ctx context.Context, q models.WorkingDaysQuery) ([]string, error) {
	return s.repo.GetWorkingDaysForMonth(ctx, q.MasterID, q.Month)
}

func (s *ScheduleService) GetForMonth(ctx context.Context, masterID, yearMonth string) (*models.MonthlySchedule, error) {
	return s.repo.GetForMonth(ctx, masterID, yearMonth)
}

func (s *ScheduleService) SaveTimeSlots(ctx context.Context, masterID string, req models.SaveTimeSlotsRequest) error {
	return s.repo.SaveTimeSlots(ctx, masterID, req.Dates, req.Slots)
}

func (s *ScheduleService) DeleteWorkingDay(ctx context.Context, masterID, date string) error {
	return s.repo.DeleteWorkingDay(ctx, masterID, date)
}

// CreateBreak создает перерыв мастера.
func (s *ScheduleService) CreateBreak(ctx context.Context, masterID string, req models.CreateMasterBreakRequest) (*models.MasterBreak, error) {
	b := &models.MasterBreak{
		MasterID:    masterID,
		BreakTypeID: req.BreakTypeID,
		StartTime:   req.StartTime.UTC(),
		EndTime:     req.EndTime.UTC(),
	}
	if err := s.repo.CreateBreak(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetBreaks возвращает перерывы мастера.
func (s *ScheduleService) GetBreaks(ctx context.Context, masterID string) ([]models.MasterBreak, error) {
	return s.repo.GetBreaks(ctx, masterID)
}

// DeleteBreak удаляет перерыв.
func (s *ScheduleService) DeleteBreak(ctx context.Context, id, masterID string) error {
	return s.repo.DeleteBreak(ctx, id, masterID)
}
