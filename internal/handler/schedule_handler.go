package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-backend/internal/models"
	"booking-backend/internal/service"
)

type ScheduleHandler struct {
	svc *service.ScheduleService
}

func NewScheduleHandler(svc *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

// GetForMonth godoc
// @Summary      Расписание на месяц
// @Tags         Schedule
// @Produce      json
// @Security     BearerAuth
// @Param        year_month  query  string  true  "Месяц YYYY-MM"
// @Success      200  {object}  models.MonthlySchedule
// @Router       /schedule [get]
func (h *ScheduleHandler) GetForMonth(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	yearMonth := c.Query("year_month")
	if yearMonth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр year_month обязателен"})
		return
	}

	schedule, err := h.svc.GetForMonth(c.Request.Context(), masterID, yearMonth)
	if err != nil {
		c.JSON(http.StatusOK, models.MonthlySchedule{
			MasterID:    masterID,
			YearMonth:   yearMonth,
			WorkingDays: map[string][]models.TimeSlot{},
		})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// SaveTimeSlots godoc
// @Summary      Сохранить рабочие слоты
// @Tags         Schedule
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  models.SaveTimeSlotsRequest  true  "Даты и слоты"
// @Success      200   {object}  models.SuccessResponse
// @Router       /schedule/slots [put]
func (h *ScheduleHandler) SaveTimeSlots(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.SaveTimeSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SaveTimeSlots(c.Request.Context(), masterID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "SaveTimeSlots")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "расписание сохранено"})
}

// DeleteWorkingDay godoc
// @Summary      Удалить рабочий день
// @Tags         Schedule
// @Produce      json
// @Security     BearerAuth
// @Param        date  query  string  true  "Дата YYYY-MM-DD"
// @Success      200   {object}  models.SuccessResponse
// @Router       /schedule/day [delete]
func (h *ScheduleHandler) DeleteWorkingDay(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр date обязателен"})
		return
	}

	if err := h.svc.DeleteWorkingDay(c.Request.Context(), masterID, date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "DeleteWorkingDay")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "рабочий день удалён"})
}

// GetAvailableSlots godoc
// @Summary      Доступные слоты для записи
// @Tags         Schedule
// @Produce      json
// @Param        master_id   query  string  true  "ID мастера"
// @Param        date        query  string  true  "Дата YYYY-MM-DD"
// @Success      200  {object}  object{slots=[]string}
// @Router       /schedule/available-slots [get]
func (h *ScheduleHandler) GetAvailableSlots(c *gin.Context) {
	var query models.AvailableSlotsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slots, err := h.svc.GetAvailableSlots(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetAvailableSlots")})
		return
	}

	if slots == nil {
		slots = []string{}
	}

	c.JSON(http.StatusOK, gin.H{"slots": slots})
}

// GetWorkingDays godoc
// @Summary      Рабочие дни месяца
// @Tags         Schedule
// @Produce      json
// @Param        master_id  query  string  true  "ID мастера"
// @Param        month      query  string  true  "Месяц YYYY-MM"
// @Success      200  {object}  object{working_days=[]string}
// @Router       /schedule/working-days [get]
func (h *ScheduleHandler) GetWorkingDays(c *gin.Context) {
	var query models.WorkingDaysQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	days, err := h.svc.GetWorkingDays(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetWorkingDays")})
		return
	}

	if days == nil {
		days = []string{}
	}

	c.JSON(http.StatusOK, gin.H{"working_days": days})
}

// CreateBreak godoc
// @Summary      Создать перерыв мастера
// @Tags         Breaks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  models.CreateMasterBreakRequest  true  "Параметры перерыва"
// @Success      201   {object}  models.MasterBreak
// @Router       /breaks [post]
func (h *ScheduleHandler) CreateBreak(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.CreateMasterBreakRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b, err := h.svc.CreateBreak(c.Request.Context(), masterID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "CreateBreak")})
		return
	}

	c.JSON(http.StatusCreated, b)
}

// GetBreaks godoc
// @Summary      Список перерывов мастера
// @Tags         Breaks
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.MasterBreak
// @Router       /breaks [get]
func (h *ScheduleHandler) GetBreaks(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	list, err := h.svc.GetBreaks(c.Request.Context(), masterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetBreaks")})
		return
	}

	if list == nil {
		list = []models.MasterBreak{}
	}

	c.JSON(http.StatusOK, list)
}

// DeleteBreak godoc
// @Summary      Удалить перерыв
// @Tags         Breaks
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID перерыва"
// @Success      200  {object}  models.SuccessResponse
// @Router       /breaks/{id} [delete]
func (h *ScheduleHandler) DeleteBreak(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	if err := h.svc.DeleteBreak(c.Request.Context(), id, masterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "перерыв удалён"})
}
