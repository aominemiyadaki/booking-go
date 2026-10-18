package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-backend/internal/models"
	"booking-backend/internal/service"
)

type ServiceHandler struct {
	svc *service.ServiceService
}

func NewServiceHandler(svc *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

// GetByMaster godoc
// @Summary      Список услуг
// @Tags         Services
// @Produce      json
// @Param        master_id  query  string  false  "ID мастера (для публичного доступа)"
// @Success      200  {array}  models.Service
// @Router       /services [get]
func (h *ServiceHandler) GetByMaster(c *gin.Context) {
	masterID := c.Query("master_id")
	if masterID == "" {
		if val, exists := c.Get("master_id"); exists {
			if idStr, ok := val.(string); ok {
				masterID = idStr
			}
		}
	}

	if masterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не указан master_id"})
		return
	}

	services, err := h.svc.GetByMaster(c.Request.Context(), masterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetByMaster services")})
		return
	}

	if services == nil {
		services = []models.Service{}
	}

	c.JSON(http.StatusOK, services)
}

// GetByID godoc
// @Summary      Детали услуги
// @Tags         Services
// @Produce      json
// @Param        id   path      string  true  "ID услуги"
// @Success      200  {object}  models.Service
// @Failure      404  {object}  models.ErrorResponse
// @Router       /services/{id} [get]
func (h *ServiceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	svc, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "услуга не найдена"})
		return
	}

	c.JSON(http.StatusOK, svc)
}

// Create godoc
// @Summary      Добавить услугу
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  models.CreateServiceRequest  true  "Данные услуги"
// @Success      201   {object}  models.Service
// @Router       /services [post]
func (h *ServiceHandler) Create(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc, err := h.svc.Create(c.Request.Context(), masterID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "Create service")})
		return
	}

	c.JSON(http.StatusCreated, svc)
}

// Update godoc
// @Summary      Обновить услугу
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                      true  "ID услуги"
// @Param        body  body  models.UpdateServiceRequest  true  "Новые данные"
// @Success      200   {object}  models.SuccessResponse
// @Router       /services/{id} [put]
func (h *ServiceHandler) Update(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	var req models.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(c.Request.Context(), id, masterID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "услуга обновлена"})
}

// Delete godoc
// @Summary      Удалить услугу
// @Tags         Services
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID услуги"
// @Success      200  {object}  models.SuccessResponse
// @Router       /services/{id} [delete]
func (h *ServiceHandler) Delete(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id, masterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "услуга удалена"})
}
