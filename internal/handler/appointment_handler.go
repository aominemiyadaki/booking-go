package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"booking-backend/internal/models"
	"booking-backend/internal/service"
)

type AppointmentHandler struct {
	svc *service.AppointmentService
}

func NewAppointmentHandler(svc *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{svc: svc}
}

// CreateMaster godoc
// @Summary      Создать запись (мастер)
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.BookAppointmentRequest  true  "Данные записи"
// @Success      201   {object}  models.Appointment
// @Failure      400   {object}  models.ErrorResponse
// @Failure      409   {object}  models.ErrorResponse
// @Router       /appointments [post]
func (h *AppointmentHandler) CreateMaster(c *gin.Context) {
	_, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// CreateClient godoc
// @Summary      Создать запись (клиент)
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        body  body      models.BookAppointmentRequest  true  "Данные записи"
// @Success      201   {object}  models.Appointment
// @Failure      400   {object}  models.ErrorResponse
// @Failure      409   {object}  models.ErrorResponse
// @Router       /client/appointments [post]
func (h *AppointmentHandler) CreateClient(c *gin.Context) {
	var req models.BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// GetByID godoc
// @Summary      Детали записи
// @Tags         Appointments
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID записи"
// @Success      200  {object}  models.Appointment
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /appointments/{id} [get]
func (h *AppointmentHandler) GetByID(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	appointment, err := h.svc.GetByID(c.Request.Context(), id, masterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// Update godoc
// @Summary      Редактировать/перенести запись
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                           true  "ID записи"
// @Param        body  body      models.BookAppointmentRequest  true  "Новые параметры записи"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Router       /appointments/{id} [put]
func (h *AppointmentHandler) Update(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	var req models.BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(c.Request.Context(), id, masterID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "запись успешно обновлена"})
}

// GetByDate godoc
// @Summary      Записи на дату
// @Tags         Appointments
// @Produce      json
// @Security     BearerAuth
// @Param        date  query     string  true  "Дата в формате YYYY-MM-DD"
// @Success      200   {array}   models.Appointment
// @Failure      400   {object}  models.ErrorResponse
// @Router       /appointments [get]
func (h *AppointmentHandler) GetByDate(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр date обязателен"})
		return
	}

	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат даты, ожидается YYYY-MM-DD"})
		return
	}

	appointments, err := h.svc.GetByDate(c.Request.Context(), masterID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetByDate")})
		return
	}

	if appointments == nil {
		appointments = []models.AppointmentResponse{}
	}

	c.JSON(http.StatusOK, appointments)
}

// GetGuestAppointments godoc
// @Summary      Мои записи (гость)
// @Tags         Appointments
// @Produce      json
// @Param        guest_token  query  string  false  "Токен гостя"
// @Param        phone        query  string  false  "Телефон гостя"
// @Success      200          {array}   models.Appointment
// @Failure      400          {object}  models.ErrorResponse
// @Router       /appointments/my [get]
func (h *AppointmentHandler) GetGuestAppointments(c *gin.Context) {
	guestToken := c.Query("guest_token")
	phone := c.Query("phone")
	if guestToken == "" && phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр guest_token или phone обязателен"})
		return
	}

	appointments, err := h.svc.GetGuestAppointments(c.Request.Context(), guestToken, phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetGuestAppointments")})
		return
	}

	if appointments == nil {
		appointments = []models.AppointmentResponse{}
	}

	c.JSON(http.StatusOK, appointments)
}

// UpdateStatus godoc
// @Summary      Изменить статус записи
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                      true  "ID записи"
// @Param        body  body      models.UpdateStatusRequest   true  "Новый статус"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Router       /appointments/{id}/status [patch]
func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusStr := req.StatusID
	if statusStr == "" {
		statusStr = req.Status
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), id, masterID, statusStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "статус обновлён"})
}

// Delete godoc
// @Summary      Удалить запись
// @Tags         Appointments
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID записи"
// @Success      200  {object}  models.SuccessResponse
// @Router       /appointments/{id} [delete]
func (h *AppointmentHandler) Delete(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id, masterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "запись удалена"})
}

// CreateBlock godoc
// @Summary      Создать блок/перерыв
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      interface{}  true  "Время блока"
// @Success      201   {object}  models.Appointment
// @Failure      409   {object}  models.ErrorResponse
// @Router       /appointments/block [post]
func (h *AppointmentHandler) CreateBlock(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	block, err := h.svc.CreateBlock(c.Request.Context(), masterID, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, block)
}
