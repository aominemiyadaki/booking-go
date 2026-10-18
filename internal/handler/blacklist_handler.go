package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-backend/internal/models"
	"booking-backend/internal/service"
)

type BlacklistHandler struct {
	svc *service.BlacklistService
}

func NewBlacklistHandler(svc *service.BlacklistService) *BlacklistHandler {
	return &BlacklistHandler{svc: svc}
}

// GetAll godoc
// @Summary      Список заблокированных (мастера)
// @Tags         Blacklist
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.BlacklistedClient
// @Router       /blacklist [get]
func (h *BlacklistHandler) GetAll(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	list, err := h.svc.GetAll(c.Request.Context(), masterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetAll blacklist")})
		return
	}

	if list == nil {
		list = []models.BlacklistedClient{}
	}

	c.JSON(http.StatusOK, list)
}

// Block godoc
// @Summary      Заблокировать клиента
// @Tags         Blacklist
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  models.BlockClientRequest  true  "Данные для блокировки"
// @Success      200   {object}  models.SuccessResponse
// @Router       /blacklist [post]
func (h *BlacklistHandler) Block(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.BlockClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Block(c.Request.Context(), masterID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "Block")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "клиент заблокирован"})
}

// Unblock godoc
// @Summary      Разблокировать клиента
// @Tags         Blacklist
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "ID записи в blacklist"
// @Success      200  {object}  models.SuccessResponse
// @Router       /blacklist/{id} [delete]
func (h *BlacklistHandler) Unblock(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	if err := h.svc.Unblock(c.Request.Context(), id, masterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "клиент разблокирован"})
}

// Check godoc
// @Summary      Проверка blacklist
// @Tags         Blacklist
// @Produce      json
// @Param        master_id    query  string  true   "ID мастера"
// @Param        phone        query  string  false  "Номер телефона"
// @Param        guest_token  query  string  false  "Guest token"
// @Success      200          {object}  object{blocked=bool}
// @Router       /blacklist/check [get]
func (h *BlacklistHandler) Check(c *gin.Context) {
	masterID := c.Query("master_id")
	if masterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр master_id обязателен"})
		return
	}

	phone := c.Query("phone")
	guestToken := c.Query("guest_token")

	blocked, err := h.svc.Check(c.Request.Context(), masterID, phone, guestToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "Check blacklist")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocked": blocked})
}
