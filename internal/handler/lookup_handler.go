package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-backend/internal/models"
	"booking-backend/internal/service"
)

type LookupHandler struct {
	svc *service.LookupService
}

func NewLookupHandler(svc *service.LookupService) *LookupHandler {
	return &LookupHandler{svc: svc}
}

// === Specializations ===

// GetSpecializations godoc
// @Summary      Список специализаций
// @Tags         Specializations
// @Produce      json
// @Success      200  {array}  models.Specialization
// @Router       /specializations [get]
func (h *LookupHandler) GetSpecializations(c *gin.Context) {
	list, err := h.svc.GetSpecializations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetSpecializations")})
		return
	}
	c.JSON(http.StatusOK, list)
}

// CreateSpecialization godoc
// @Summary      Добавить специализацию
// @Tags         Specializations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  object{name=string}  true  "Название специализации"
// @Success      201   {object}  models.Specialization
// @Router       /specializations [post]
func (h *LookupHandler) CreateSpecialization(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.CreateSpecialization(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "CreateSpecialization")})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateSpecialization godoc
// @Summary      Обновить специализацию
// @Tags         Specializations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                  true  "ID специализации"
// @Param        body  body  object{name=string}  true  "Новое название"
// @Success      200   {object}  models.SuccessResponse
// @Router       /specializations/{id} [put]
func (h *LookupHandler) UpdateSpecialization(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateSpecialization(c.Request.Context(), id, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "UpdateSpecialization")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "специализация обновлена"})
}

// DeleteSpecialization godoc
// @Summary      Удалить специализацию
// @Tags         Specializations
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID специализации"
// @Success      200  {object}  models.SuccessResponse
// @Router       /specializations/{id} [delete]
func (h *LookupHandler) DeleteSpecialization(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSpecialization(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "DeleteSpecialization")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "специализация удалена"})
}

// === Currencies ===

// GetCurrencies godoc
// @Summary      Список валют
// @Tags         Currencies
// @Produce      json
// @Success      200  {array}  models.Currency
// @Router       /currencies [get]
func (h *LookupHandler) GetCurrencies(c *gin.Context) {
	list, err := h.svc.GetCurrencies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetCurrencies")})
		return
	}
	c.JSON(http.StatusOK, list)
}

// CreateCurrency godoc
// @Summary      Добавить валюту
// @Tags         Currencies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  object{name=string,iso_code=string}  true  "Параметры валюты"
// @Success      201   {object}  models.Currency
// @Router       /currencies [post]
func (h *LookupHandler) CreateCurrency(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		ISOCode string `json:"iso_code" binding:"required,len=3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.CreateCurrency(c.Request.Context(), req.Name, req.ISOCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "CreateCurrency")})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateCurrency godoc
// @Summary      Обновить валюту
// @Tags         Currencies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                                  true  "ID валюты"
// @Param        body  body  object{name=string,iso_code=string}  true  "Новые параметры"
// @Success      200   {object}  models.SuccessResponse
// @Router       /currencies/{id} [put]
func (h *LookupHandler) UpdateCurrency(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name    string `json:"name" binding:"required"`
		ISOCode string `json:"iso_code" binding:"required,len=3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateCurrency(c.Request.Context(), id, req.Name, req.ISOCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "UpdateCurrency")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "валюта обновлена"})
}

// DeleteCurrency godoc
// @Summary      Удалить валюту
// @Tags         Currencies
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID валюты"
// @Success      200  {object}  models.SuccessResponse
// @Router       /currencies/{id} [delete]
func (h *LookupHandler) DeleteCurrency(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCurrency(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "DeleteCurrency")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "валюта удалена"})
}

// === AppointmentStatuses ===

// GetAppointmentStatuses godoc
// @Summary      Список статусов записей
// @Tags         Appointments
// @Produce      json
// @Success      200  {array}  models.AppointmentStatus
// @Router       /appointment-statuses [get]
func (h *LookupHandler) GetAppointmentStatuses(c *gin.Context) {
	list, err := h.svc.GetAppointmentStatuses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetAppointmentStatuses")})
		return
	}
	c.JSON(http.StatusOK, list)
}

// CreateAppointmentStatus godoc
// @Summary      Добавить статус записи
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  object{name=string}  true  "Название статуса"
// @Success      201   {object}  models.AppointmentStatus
// @Router       /appointment-statuses [post]
func (h *LookupHandler) CreateAppointmentStatus(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.CreateAppointmentStatus(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "CreateAppointmentStatus")})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateAppointmentStatus godoc
// @Summary      Обновить статус записи
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                  true  "ID статуса"
// @Param        body  body  object{name=string}  true  "Новое название"
// @Success      200   {object}  models.SuccessResponse
// @Router       /appointment-statuses/{id} [put]
func (h *LookupHandler) UpdateAppointmentStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateAppointmentStatus(c.Request.Context(), id, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "UpdateAppointmentStatus")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "статус записи обновлен"})
}

// DeleteAppointmentStatus godoc
// @Summary      Удалить статус записи
// @Tags         Appointments
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID статуса"
// @Success      200  {object}  models.SuccessResponse
// @Router       /appointment-statuses/{id} [delete]
func (h *LookupHandler) DeleteAppointmentStatus(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAppointmentStatus(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "DeleteAppointmentStatus")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "статус записи удален"})
}

// === BreakTypes ===

// GetBreakTypes godoc
// @Summary      Список типов перерывов
// @Tags         Breaks
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.BreakType
// @Router       /break-types [get]
func (h *LookupHandler) GetBreakTypes(c *gin.Context) {
	list, err := h.svc.GetBreakTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetBreakTypes")})
		return
	}
	c.JSON(http.StatusOK, list)
}

// CreateBreakType godoc
// @Summary      Добавить тип перерыва
// @Tags         Breaks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  object{name=string}  true  "Название типа"
// @Success      201   {object}  models.BreakType
// @Router       /break-types [post]
func (h *LookupHandler) CreateBreakType(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.CreateBreakType(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "CreateBreakType")})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateBreakType godoc
// @Summary      Обновить тип перерыва
// @Tags         Breaks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                  true  "ID типа"
// @Param        body  body  object{name=string}  true  "Новое название"
// @Success      200   {object}  models.SuccessResponse
// @Router       /break-types/{id} [put]
func (h *LookupHandler) UpdateBreakType(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateBreakType(c.Request.Context(), id, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "UpdateBreakType")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "тип перерыва обновлен"})
}

// DeleteBreakType godoc
// @Summary      Удалить тип перерыва
// @Tags         Breaks
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID типа"
// @Success      200  {object}  models.SuccessResponse
// @Router       /break-types/{id} [delete]
func (h *LookupHandler) DeleteBreakType(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteBreakType(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "DeleteBreakType")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "тип перерыва удален"})
}
