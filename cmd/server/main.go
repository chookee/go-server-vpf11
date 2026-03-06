// Package main - точка входа приложения.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zerocode/users-api/internal/config"
	"github.com/zerocode/users-api/internal/database"
	"github.com/zerocode/users-api/internal/handlers"
	appMiddleware "github.com/zerocode/users-api/internal/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	logger, err := initLogger(cfg.Log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации логгера: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Делаем логгер глобально доступным
	zap.ReplaceGlobals(logger)

	logger.Info("starting Users API server",
		zap.String("host", cfg.Server.Host),
		zap.String("port", cfg.Server.Port),
	)

	// Инициализация базы данных
	db, err := database.New(cfg.Database.Path, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, logger)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Init(); err != nil {
		logger.Fatal("failed to initialize database schema", zap.Error(err))
	}

	// Создание обработчиков
	h := handlers.New(db.DB, logger)

	// Настройка роутера
	r := chi.NewRouter()

	// Применение middleware
	r.Use(appMiddleware.Chain(logger))
	
	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting (100 запросов в секунду) - используем встроенный Throttle
	r.Use(middleware.Throttle(100))

	// Валидация Content-Type для POST/PUT
	r.Use(appMiddleware.ContentTypeValidator("application/json"))

	// Ограничение размера тела запроса
	r.Use(appMiddleware.RequestLimit(cfg.Server.MaxRequestBody))

	// Health check
	r.Get("/health", h.HealthCheck)

	// API routes
	r.Group(func(r chi.Router) {
		r.Get("/users", h.GetUsers)
		r.Post("/add-user", h.AddUser)
		r.Get("/user/{id}", h.GetUser)
		r.Put("/user/{id}", h.UpdateUser)
		r.Delete("/user/{id}", h.DeleteUser)
		r.Get("/stats/active", h.GetStats)
	})

	// Создание HTTP сервера
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Канал для обработки сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		logger.Info("server started", zap.String("address", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	logger.Info("shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}

	// Закрываем соединение с БД
	if err := db.Close(); err != nil {
		logger.Error("failed to close database connection", zap.Error(err))
	}

	logger.Info("server stopped gracefully")
}

// initLogger инициализирует zap logger.
func initLogger(cfg config.LogConfig) (*zap.Logger, error) {
	var loggerConfig zap.Config
	
	switch cfg.Level {
	case "debug":
		loggerConfig = zap.NewDevelopmentConfig()
	case "production":
		loggerConfig = zap.NewProductionConfig()
	default:
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Настраиваем encoder для читаемости в development
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if cfg.File != "" {
		loggerConfig.OutputPaths = []string{cfg.File}
		loggerConfig.ErrorOutputPaths = []string{cfg.File}
	}

	return loggerConfig.Build()
}
