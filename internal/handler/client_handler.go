package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-backend/internal/models"
	"booking-backend/internal/service"
)

type ClientHandler struct {
	svc *service.ClientService
}

func NewClientHandler(svc *service.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

// GetByMaster godoc
// @Summary      Список клиентов (с пагинацией)
// @Tags         Clients
// @Produce      json
// @Security     BearerAuth
// @Param        page       query  int  false  "Номер страницы"  default(1)
// @Param        page_size  query  int  false  "Размер страницы"  default(50)
// @Success      200  {object}  models.PaginatedResponse
// @Router       /clients [get]
func (h *ClientHandler) GetByMaster(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var page models.PaginationRequest
	_ = c.ShouldBindQuery(&page)
	page.Normalize()

	clients, total, err := h.svc.GetByMaster(c.Request.Context(), masterID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetByMaster clients")})
		return
	}

	if clients == nil {
		clients = []models.Client{}
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       clients,
		Total:      total,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: models.TotalPages(total, page.PageSize),
	})
}

// GetByID godoc
// @Summary      Детали клиента
// @Tags         Clients
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "ID клиента"
// @Success      200  {object}  models.Client
// @Failure      404  {object}  models.ErrorResponse
// @Router       /clients/{id} [get]
func (h *ClientHandler) GetByID(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	client, err := h.svc.GetByID(c.Request.Context(), id, masterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "клиент не найден"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// Create godoc
// @Summary      Добавить клиента
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  models.CreateClientRequest  true  "Данные клиента"
// @Success      201   {object}  models.Client
// @Router       /clients [post]
func (h *ClientHandler) Create(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.svc.Create(c.Request.Context(), masterID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "Create client")})
		return
	}

	c.JSON(http.StatusCreated, client)
}

// Update godoc
// @Summary      Обновить клиента
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                     true  "ID клиента"
// @Param        body  body  models.UpdateClientRequest  true  "Новые данные"
// @Success      200   {object}  models.SuccessResponse
// @Router       /clients/{id} [put]
func (h *ClientHandler) Update(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	var req models.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(c.Request.Context(), id, masterID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "клиент обновлён"})
}

// Delete godoc
// @Summary      Удалить клиента
// @Tags         Clients
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID клиента"
// @Success      200  {object}  models.SuccessResponse
// @Router       /clients/{id} [delete]
func (h *ClientHandler) Delete(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id, masterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "клиент удалён"})
}

// UpdateNotes godoc
// @Summary      Обновить заметки
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                    true  "ID клиента"
// @Param        body  body  models.UpdateNotesRequest  true  "Текст заметки"
// @Success      200   {object}  models.SuccessResponse
// @Router       /clients/{id}/notes [patch]
func (h *ClientHandler) UpdateNotes(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	var req models.UpdateNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateNotes(c.Request.Context(), id, masterID, req.Notes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "заметки обновлены"})
}

// ToggleBlock godoc
// @Summary      Заблокировать/разблокировать клиента
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                     true  "ID клиента"
// @Param        body  body  models.ToggleBlockRequest   true  "Статус блокировки"
// @Success      200   {object}  models.SuccessResponse
// @Router       /clients/{id}/block [patch]
func (h *ClientHandler) ToggleBlock(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	var req models.ToggleBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ToggleBlock(c.Request.Context(), id, masterID, req.Block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := "разблокирован"
	if req.Block {
		action = "заблокирован"
	}
	c.JSON(http.StatusOK, gin.H{"message": "клиент " + action})
}

// GetHistory godoc
// @Summary      История визитов (с пагинацией)
// @Tags         Clients
// @Produce      json
// @Security     BearerAuth
// @Param        id         path   string  true   "ID клиента"
// @Param        page       query  int     false  "Номер страницы"  default(1)
// @Param        page_size  query  int     false  "Размер страницы"  default(50)
// @Success      200  {object}  models.PaginatedResponse
// @Router       /clients/{id}/history [get]
func (h *ClientHandler) GetHistory(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	var page models.PaginationRequest
	_ = c.ShouldBindQuery(&page)
	page.Normalize()

	history, total, err := h.svc.GetHistory(c.Request.Context(), id, masterID, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if history == nil {
		history = []models.AppointmentResponse{}
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       history,
		Total:      total,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: models.TotalPages(total, page.PageSize),
	})
}

// FindOrCreate godoc
// @Summary      Найти или создать клиента
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        body  body  models.FindOrCreateClientRequest  true  "Данные клиента"
// @Success      200   {object}  models.FindOrCreateClientResponse
// @Router       /clients/find-or-create [post]
func (h *ClientHandler) FindOrCreate(c *gin.Context) {
	var req models.FindOrCreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.svc.FindOrCreate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "FindOrCreate client")})
		return
	}

	c.JSON(http.StatusOK, models.FindOrCreateClientResponse{
		ClientID: client.ID,
	})
}
