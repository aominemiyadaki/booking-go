package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"booking-backend/internal/models"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) GetWorkingSlotsForDate(ctx context.Context, masterID, date string) ([]models.TimeSlot, error) {
	var day models.WorkingDay
	err := r.db.WithContext(ctx).
		Where("master_id = ? AND date = ? AND visibility = true", masterID, date).
		First(&day).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Выходной
		}
		return nil, err
	}

	var slots []models.WorkingSlot
	err = r.db.WithContext(ctx).
		Where("working_day_id = ? AND visibility = true", day.ID).
		Order("start_time").
		Find(&slots).Error
	if err != nil {
		return nil, err
	}

	var result []models.TimeSlot
	for _, s := range slots {
		result = append(result, models.TimeSlot{
			Start: s.StartTime,
			End:   s.EndTime,
		})
	}
	return result, nil
}

func (r *ScheduleRepository) SaveTimeSlots(ctx context.Context, masterID string, dates []string, slots []models.TimeSlot) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, dateStr := range dates {
			// Парсим строку даты в формат time.Time
			parsedDate, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
			if err != nil {
				return err // Возвращаем ошибку, если дата пришла в неверном формате
			}

			var day models.WorkingDay
			// Ищем день для этого мастера на указанную дату
			err = tx.Where("master_id = ? AND date = ?", masterID, parsedDate).First(&day).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Если дня нет — создаем новый
					day = models.WorkingDay{
						MasterID:   masterID,
						Date:       parsedDate,
						Visibility: true,
					}
					if err := tx.Create(&day).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				// Если день найден, обновляем видимость (вдруг он был удален)
				if !day.Visibility {
					tx.Model(&day).Update("visibility", true)
				}
				// Удаляем (или скрываем) все старые слоты для этого дня перед сохранением новых
				if err := tx.Where("working_day_id = ?", day.ID).Delete(&models.WorkingSlot{}).Error; err != nil {
					return err
				}
			}

			// Создаем новые слоты для этого дня
			for _, slot := range slots {
				newSlot := models.WorkingSlot{
					WorkingDayID: day.ID,
					StartTime:    slot.Start,
					EndTime:      slot.End,
				}
				if err := tx.Create(&newSlot).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *ScheduleRepository) DeleteWorkingDay(ctx context.Context, masterID, date string) error {
	parsedDate, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return err
	}

	// Soft delete (скрываем день через visibility)
	result := r.db.WithContext(ctx).
		Model(&models.WorkingDay{}).
		Where("master_id = ? AND date = ?", masterID, parsedDate).
		Update("visibility", false)
	
	return result.Error
}

func getMonthDateRange(yearMonth string) (time.Time, time.Time, error) {
	t, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	return start, end, nil
}

func (r *ScheduleRepository) GetWorkingDaysForMonth(ctx context.Context, masterID, yearMonth string) ([]string, error) {
	start, end, err := getMonthDateRange(yearMonth)
	if err != nil {
		return nil, err
	}

	var dates []time.Time
	err = r.db.WithContext(ctx).
		Model(&models.WorkingDay{}).
		Where("master_id = ? AND visibility = true AND date >= ? AND date < ?", masterID, start, end).
		Order("date ASC").
		Pluck("date", &dates).Error
	
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(dates))
	for _, d := range dates {
		result = append(result, d.Format("2006-01-02"))
	}
	
	return result, nil
}

func (r *ScheduleRepository) GetForMonth(ctx context.Context, masterID, yearMonth string) (*models.MonthlySchedule, error) {
	start, end, err := getMonthDateRange(yearMonth)
	if err != nil {
		return nil, err
	}

	var days []models.WorkingDay
	err = r.db.WithContext(ctx).
		Where("master_id = ? AND visibility = true AND date >= ? AND date < ?", masterID, start, end).
		Order("date ASC").
		Find(&days).Error
	if err != nil {
		return nil, err
	}

	result := &models.MonthlySchedule{
		MasterID:    masterID,
		YearMonth:   yearMonth,
		WorkingDays: make(map[string][]models.TimeSlot),
	}

	if len(days) == 0 {
		return result, nil
	}

	dayIDs := make([]string, 0, len(days))
	dayMap := make(map[string]string, len(days))
	for _, day := range days {
		dayIDs = append(dayIDs, day.ID)
		dayMap[day.ID] = day.Date.Format("2006-01-02")
	}

	var slots []models.WorkingSlot
	err = r.db.WithContext(ctx).
		Where("working_day_id IN ? AND visibility = true", dayIDs).
		Order("start_time ASC").
		Find(&slots).Error
	if err != nil {
		return nil, err
	}

	for _, s := range slots {
		dateStr, ok := dayMap[s.WorkingDayID]
		if !ok {
			continue
		}
		result.WorkingDays[dateStr] = append(result.WorkingDays[dateStr], models.TimeSlot{
			Start: s.StartTime,
			End:   s.EndTime,
		})
	}

	return result, nil
}

// CreateBreak создает перерыв мастера в БД.
func (r *ScheduleRepository) CreateBreak(ctx context.Context, b *models.MasterBreak) error {
	return r.db.WithContext(ctx).Create(b).Error
}

// GetBreaks возвращает список активных перерывов мастера.
func (r *ScheduleRepository) GetBreaks(ctx context.Context, masterID string) ([]models.MasterBreak, error) {
	var list []models.MasterBreak
	err := r.db.WithContext(ctx).
		Where("master_id = ? AND visibility = true", masterID).
		Order("start_time").
		Find(&list).Error
	return list, err
}

// DeleteBreak выполняет мягкое удаление перерыва.
func (r *ScheduleRepository) DeleteBreak(ctx context.Context, id, masterID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.MasterBreak{}).
		Where("id = ? AND master_id = ?", id, masterID).
		Update("visibility", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("перерыв не найден или не принадлежит мастеру")
	}
	return nil
}
