package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"booking-backend/internal/models"
	"booking-backend/internal/repository"
)

type MasterService struct {
	repo      repository.MasterRepositoryInterface
	jwtSecret string
}

func NewMasterService(repo repository.MasterRepositoryInterface, jwtSecret string) *MasterService {
	return &MasterService{repo: repo, jwtSecret: jwtSecret}
}

func (s *MasterService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	reqID, _ := ctx.Value("request_id").(string)

	existing, _ := s.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		slog.Warn("Попытка регистрации с существующим email", "email", req.Email, "request_id", reqID)
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// SECURITY: Проверяем сложность пароля
	if len(req.Password) < 8 {
		return nil, errors.New("пароль должен содержать минимум 8 символов")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Ошибка хэширования пароля", "error", err, "request_id", reqID)
		return nil, errors.New("ошибка сервера")
	}

	master := &models.Master{
		Email:            req.Email,
		PasswordHash:     string(hash),
		FullName:         req.FullName,
		PhoneNumber:      req.PhoneNumber,
		SpecializationID: req.SpecializationID,
	}

	if err := s.repo.Create(ctx, master); err != nil {
		slog.Error("Не удалось сохранить мастера в БД", "error", err, "request_id", reqID)
		return nil, err
	}

	// Перечитываем из базы для загрузки связи Specialization
	if loaded, err := s.repo.GetByID(ctx, master.ID); err == nil {
		master = loaded
	}

	token, err := s.generateToken(master.ID)
	if err != nil {
		slog.Error("Ошибка генерации токена", "error", err, "request_id", reqID)
		return nil, errors.New("ошибка сервера")
	}

	slog.Info("Мастер успешно зарегистрирован", "master_id", master.ID, "request_id", reqID)

	return &models.AuthResponse{
		Token: token,
		Master: models.MasterPublic{
			ID:               master.ID,
			FullName:         master.FullName,
			PhoneNumber:      master.PhoneNumber,
			SpecializationID: master.SpecializationID,
			Specialization:   master.SpecializationName,
			AvatarURL:        master.AvatarURL,
		},
	}, nil
}

func (s *MasterService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	reqID, _ := ctx.Value("request_id").(string)

	master, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		slog.Warn("Логин не удался: email не найден", "email", req.Email, "request_id", reqID)
		return nil, errors.New("неверный email или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(master.PasswordHash), []byte(req.Password)); err != nil {
		slog.Warn("Логин не удался: неверный пароль", "email", req.Email, "request_id", reqID)
		return nil, errors.New("неверный email или пароль")
	}

	token, err := s.generateToken(master.ID)
	if err != nil {
		slog.Error("Ошибка генерации токена", "error", err, "request_id", reqID)
		return nil, errors.New("ошибка сервера")
	}

	slog.Info("Мастер успешно вошёл в систему", "master_id", master.ID, "request_id", reqID)

	return &models.AuthResponse{
		Token: token,
		Master: models.MasterPublic{
			ID:               master.ID,
			FullName:         master.FullName,
			PhoneNumber:      master.PhoneNumber,
			SpecializationID: master.SpecializationID,
			Specialization:   master.SpecializationName,
			AvatarURL:        master.AvatarURL,
		},
	}, nil
}

func (s *MasterService) GetProfile(ctx context.Context, masterID string) (*models.MasterPublic, error) {
	master, err := s.repo.GetByID(ctx, masterID)
	if err != nil {
		return nil, errors.New("мастер не найден")
	}

	return &models.MasterPublic{
		ID:               master.ID,
		FullName:         master.FullName,
		PhoneNumber:      master.PhoneNumber,
		SpecializationID: master.SpecializationID,
		Specialization:   master.SpecializationName,
		AvatarURL:        master.AvatarURL,
	}, nil
}

func (s *MasterService) GetMasters(ctx context.Context) ([]models.MasterPublic, error) {
	return s.repo.GetAll(ctx)
}

// UpdateProfile обновляет данные профиля мастера.
func (s *MasterService) UpdateProfile(ctx context.Context, masterID string, req models.UpdateMasterProfileRequest) error {
	reqID, _ := ctx.Value("request_id").(string)

	if req.ID != masterID {
		slog.Warn("Попытка изменения чужого профиля мастера", "req_id", req.ID, "master_id", masterID, "request_id", reqID)
		return errors.New("доступ запрещён: вы не можете редактировать профиль другого мастера")
	}

	master, err := s.repo.GetByID(ctx, masterID)
	if err != nil {
		return err
	}

	// Перезаписываем поля
	master.FullName = req.FullName
	master.SpecializationID = req.SpecializationID
	master.PhoneNumber = req.PhoneNumber

	if err := s.repo.Update(ctx, master); err != nil {
		slog.Error("Не удалось обновить профиль мастера в БД", "error", err, "request_id", reqID)
		return err
	}

	slog.Info("Профиль мастера успешно обновлен", "master_id", masterID, "request_id", reqID)
	return nil
}

func (s *MasterService) UpdateSessionToken(ctx context.Context, masterID, token string) error {
	return s.repo.UpdateSessionToken(ctx, masterID, token)
}

func (s *MasterService) generateToken(masterID string) (string, error) {
	claims := jwt.MapClaims{
		"master_id": masterID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
