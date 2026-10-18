package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-backend/internal/models"
	"booking-backend/internal/service"
)

type MasterHandler struct {
	svc *service.MasterService
}

func NewMasterHandler(svc *service.MasterService) *MasterHandler {
	return &MasterHandler{svc: svc}
}

// Register godoc
// @Summary      Регистрация мастера
// @Description  Создаёт нового мастера и возвращает JWT токен
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "Данные для регистрации"
// @Success      201   {object}  models.AuthResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      409   {object}  models.ErrorResponse  "Email уже занят"
// @Router       /master/register [post]
func (h *MasterHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// SECURITY: Выставляем HttpOnly куку для защиты от XSS
	c.SetCookie("master_jwt", resp.Token, 86400*30, "/", "", false, true)

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary      Логин мастера
// @Description  Аутентификация по email/password, возвращает JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Email и пароль"
// @Success      200   {object}  models.AuthResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse  "Неверные credentials"
// @Router       /master/login [post]
func (h *MasterHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// SECURITY: Выставляем HttpOnly куку для защиты от XSS
	c.SetCookie("master_jwt", resp.Token, 86400*30, "/", "", false, true)

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary      Выход из системы
// @Description  Очищает HttpOnly куку авторизации
// @Tags         Auth
// @Produce      json
// @Success      200   {object}  models.SuccessResponse
// @Router       /master/logout [post]
func (h *MasterHandler) Logout(c *gin.Context) {
	c.SetCookie("master_jwt", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "успешный выход"})
}

// GetProfile godoc
// @Summary      Профиль мастера
// @Description  Возвращает профиль текущего авторизованного мастера
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.MasterPublic
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /master/profile [get]
func (h *MasterHandler) GetProfile(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	profile, err := h.svc.GetProfile(c.Request.Context(), masterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary      Обновить профиль мастера
// @Description  Полностью обновляет информацию в профиле мастера (имя, специализация, телефон)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.UpdateMasterProfileRequest  true  "Новые данные профиля"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Router       /master/profile [put]
func (h *MasterHandler) UpdateProfile(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.UpdateMasterProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateProfile(c.Request.Context(), masterID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "профиль успешно обновлен"})
}

// GetMasters godoc
// @Summary      Список мастеров
// @Tags         Auth
// @Produce      json
// @Success      200  {array}  models.MasterPublic
// @Router       /masters [get]
func (h *MasterHandler) GetMasters(c *gin.Context) {
	masters, err := h.svc.GetMasters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "GetMasters")})
		return
	}
	c.JSON(http.StatusOK, masters)
}

// UpdateSessionToken godoc
// @Summary      Обновить push-токен сессии
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.UpdateSessionTokenRequest  true  "Новый push-токен"
// @Success      200   {object}  models.SuccessResponse
// @Router       /master/session-token [put]
func (h *MasterHandler) UpdateSessionToken(c *gin.Context) {
	masterID, ok := requireMasterID(c)
	if !ok {
		return
	}

	var req models.UpdateSessionTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateSessionToken(c.Request.Context(), masterID, req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.InternalError(err, "UpdateSessionToken")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "токен сессии обновлен"})
}
