package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "booking-backend/docs"
	"booking-backend/internal/config"
	"booking-backend/internal/database"
	"booking-backend/internal/handler"
	"booking-backend/internal/logger"
	"booking-backend/internal/repository"
	"booking-backend/internal/service"
)

// @title           Booking App API
// @version         2.0
// @description     API Server for Booking Application — GORM, Kibana logging, tracing.
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Устанавливаем часовой пояс по умолчанию в Europe/Moscow
	loc, err := time.LoadLocation("Europe/Moscow")
	if err == nil {
		time.Local = loc
	}

	// 1. Инициализируем наш JSON-логгер для Kibana
	slogLogger := logger.New()

	// Загружаем переменные из .env
	_ = godotenv.Load()
	cfg := config.Load()

	// SECURITY: Проверяем, что JWT_SECRET изменён с дефолтного
	if cfg.JWTSecret == "super-secret-key-change-me" {
		if gin.Mode() == gin.ReleaseMode {
			log.Fatal("КРИТИЧНО: Установите переменную окружения JWT_SECRET! Дефолтное значение запрещено в production.")
		}
		slogLogger.Warn("⚠️  JWT_SECRET использует дефолтное значение — измените для production!")
	}

	// 2. Подключаемся к базе данных через GORM
	db, err := database.New(cfg)
	if err != nil {
		slogLogger.Error("Критическая ошибка старта", "error", err)
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	// 3. Инициализируем слои (Репозитории -> Сервисы -> Хэндлеры)
	masterRepo := repository.NewMasterRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	clientRepo := repository.NewClientRepository(db)
	blacklistRepo := repository.NewBlacklistRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	lookupRepo := repository.NewLookupRepository(db)

	masterSvc := service.NewMasterService(masterRepo, cfg.JWTSecret)
	appointmentSvc := service.NewAppointmentService(appointmentRepo, clientRepo, blacklistRepo, serviceRepo, lookupRepo)
	serviceSvc := service.NewServiceService(serviceRepo)
	clientSvc := service.NewClientService(clientRepo, appointmentRepo)
	blacklistSvc := service.NewBlacklistService(blacklistRepo)
	scheduleSvc := service.NewScheduleService(scheduleRepo, serviceRepo, appointmentRepo)
	lookupSvc := service.NewLookupService(lookupRepo)

	masterH := handler.NewMasterHandler(masterSvc)
	appointmentH := handler.NewAppointmentHandler(appointmentSvc)
	serviceH := handler.NewServiceHandler(serviceSvc)
	clientH := handler.NewClientHandler(clientSvc)
	blacklistH := handler.NewBlacklistHandler(blacklistSvc)
	scheduleH := handler.NewScheduleHandler(scheduleSvc)
	lookupH := handler.NewLookupHandler(lookupSvc)

	// 4. Настраиваем роутер (Gin)
	router := gin.New()

	// 5. ВАЖНО: Подключаем Middleware по порядку!

	router.Use(handler.CORSMiddleware(cfg.CORSOrigins))

	// SECURITY: Лимит размера тела запроса (1 MB) — защита от OOM атак
	router.Use(handler.MaxBodySizeMiddleware(1 << 20))

	// SECURITY: Заголовки безопасности OWASP
	router.Use(handler.SecurityHeadersMiddleware())

	// Сначала ловим паники, чтобы сервер не упал
	router.Use(handler.RecoveryMiddleware(slogLogger))

	// Генерируем request_id (трейс) для КАЖДОГО запроса
	router.Use(handler.RequestIDMiddleware())

	// Пишем лог с request_id в консоль (оттуда он уйдёт в Kibana)
	router.Use(logger.LoggingMiddleware(slogLogger))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// SECURITY: Swagger UI доступен только в dev-режиме
	if gin.Mode() != gin.ReleaseMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DocExpansion("none")))
	}

	api := router.Group("/api")
	{
		// SECURITY: RateLimiter (10 попыток в минуту на IP для защиты от брутфорса)
		authLimiter := handler.RateLimiterMiddleware(10)

		// Публичные роуты
		api.POST("/master/register", authLimiter, masterH.Register)
		api.POST("/master/login", authLimiter, masterH.Login)
		api.POST("/master/logout", masterH.Logout)
		api.GET("/masters", masterH.GetMasters)

		// Публичные роуты для расписания и блэклиста
		api.GET("/schedule/available-slots", scheduleH.GetAvailableSlots)
		api.GET("/schedule/working-days", scheduleH.GetWorkingDays)
		api.GET("/blacklist/check", blacklistH.Check)
		api.GET("/services", serviceH.GetByMaster)
		api.GET("/services/:id", serviceH.GetByID)
		api.POST("/client/appointments", appointmentH.CreateClient)
		api.POST("/client/find-or-create", clientH.FindOrCreate)
		api.GET("/appointments/my", appointmentH.GetGuestAppointments)

		// Публичные справочники
		api.GET("/specializations", lookupH.GetSpecializations)
		api.GET("/currencies", lookupH.GetCurrencies)
		api.GET("/appointment-statuses", lookupH.GetAppointmentStatuses)

		// Защищённые роуты (требуют токен)
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware(cfg.JWTSecret))
		{
			// Masters
			protected.GET("/master/profile", masterH.GetProfile)
			protected.PUT("/master/profile", masterH.UpdateProfile)
			protected.PUT("/master/session-token", masterH.UpdateSessionToken)

			// Appointments
			protected.POST("/appointments", appointmentH.CreateMaster)
			protected.GET("/appointments", appointmentH.GetByDate)
			protected.GET("/appointments/:id", appointmentH.GetByID)
			protected.PUT("/appointments/:id", appointmentH.Update)
			protected.PATCH("/appointments/:id/status", appointmentH.UpdateStatus)
			protected.DELETE("/appointments/:id", appointmentH.Delete)
			protected.POST("/appointments/block", appointmentH.CreateBlock)

			// Services
			protected.POST("/services", serviceH.Create)
			protected.PUT("/services/:id", serviceH.Update)
			protected.DELETE("/services/:id", serviceH.Delete)

			// Clients
			protected.GET("/clients", clientH.GetByMaster)
			protected.GET("/clients/:id", clientH.GetByID)
			protected.POST("/clients", clientH.Create)
			protected.PUT("/clients/:id", clientH.Update)
			protected.DELETE("/clients/:id", clientH.Delete)
			protected.PATCH("/clients/:id/notes", clientH.UpdateNotes)
			protected.PATCH("/clients/:id/block", clientH.ToggleBlock)
			protected.GET("/clients/:id/history", clientH.GetHistory)

			// Blacklist
			protected.GET("/blacklist", blacklistH.GetAll)
			protected.POST("/blacklist", blacklistH.Block)
			protected.DELETE("/blacklist/:id", blacklistH.Unblock)

			// Schedule
			protected.GET("/schedule", scheduleH.GetForMonth)
			protected.PUT("/schedule/slots", scheduleH.SaveTimeSlots)
			protected.DELETE("/schedule/day", scheduleH.DeleteWorkingDay)

			// Breaks
			protected.POST("/breaks", scheduleH.CreateBreak)
			protected.GET("/breaks", scheduleH.GetBreaks)
			protected.DELETE("/breaks/:id", scheduleH.DeleteBreak)

			// Lookups CRUD (Protected)
			protected.GET("/break-types", lookupH.GetBreakTypes)
			protected.POST("/break-types", lookupH.CreateBreakType)
			protected.PUT("/break-types/:id", lookupH.UpdateBreakType)
			protected.DELETE("/break-types/:id", lookupH.DeleteBreakType)

			protected.POST("/specializations", lookupH.CreateSpecialization)
			protected.PUT("/specializations/:id", lookupH.UpdateSpecialization)
			protected.DELETE("/specializations/:id", lookupH.DeleteSpecialization)

			protected.POST("/currencies", lookupH.CreateCurrency)
			protected.PUT("/currencies/:id", lookupH.UpdateCurrency)
			protected.DELETE("/currencies/:id", lookupH.DeleteCurrency)

			protected.POST("/appointment-statuses", lookupH.CreateAppointmentStatus)
			protected.PUT("/appointment-statuses/:id", lookupH.UpdateAppointmentStatus)
			protected.DELETE("/appointment-statuses/:id", lookupH.DeleteAppointmentStatus)
		}
	}

	// 6. Запускаем сервер
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slogLogger.Info("Сервер запущен", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slogLogger.Error("Ошибка сервера", "error", err)
		}
	}()

	// Graceful shutdown (плавное выключение)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slogLogger.Info("Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slogLogger.Error("Ошибка при выключении", "error", err)
	}
	slogLogger.Info("Сервер остановлен")
}
