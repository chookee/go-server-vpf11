// Package config предоставляет конфигурацию приложения.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Config содержит всю конфигурацию приложения.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// ServerConfig содержит настройки HTTP сервера.
type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxRequestBody  int64
}

// DatabaseConfig содержит настройки базы данных.
type DatabaseConfig struct {
	Path         string
	MaxOpenConns int
	MaxIdleConns int
}

// LogConfig содержит настройки логирования.
type LogConfig struct {
	Level string
	File  string
}

// Load загружает конфигурацию из переменных окружения или использует значения по умолчанию.
func Load() (*Config, error) {
	// Получаем директорию проекта (для разработки)
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Путь к БД относительно рабочей директории
	dbPath := getEnv("DB_PATH", filepath.Join(wd, "test.db"))

	// Парсим таймауты
	readTimeout, err := parseDuration(getEnv("READ_TIMEOUT", "15s"))
	if err != nil {
		return nil, fmt.Errorf("invalid READ_TIMEOUT: %w", err)
	}
	writeTimeout, err := parseDuration(getEnv("WRITE_TIMEOUT", "15s"))
	if err != nil {
		return nil, fmt.Errorf("invalid WRITE_TIMEOUT: %w", err)
	}
	idleTimeout, err := parseDuration(getEnv("IDLE_TIMEOUT", "60s"))
	if err != nil {
		return nil, fmt.Errorf("invalid IDLE_TIMEOUT: %w", err)
	}
	shutdownTimeout, err := parseDuration(getEnv("SHUTDOWN_TIMEOUT", "30s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}

	// Максимальный размер тела запроса (1MB по умолчанию)
	maxRequestBody := int64(1024 * 1024)
	if size := getEnv("MAX_REQUEST_BODY", ""); size != "" {
		if parsed, err := strconv.ParseInt(size, 10, 64); err == nil {
			maxRequestBody = parsed
		}
	}

	// Настройки пула соединений
	maxOpenConns := 25
	if conns := getEnv("DB_MAX_OPEN_CONNS", ""); conns != "" {
		if parsed, err := strconv.Atoi(conns); err == nil && parsed > 0 {
			maxOpenConns = parsed
		}
	}
	maxIdleConns := 5
	if conns := getEnv("DB_MAX_IDLE_CONNS", ""); conns != "" {
		if parsed, err := strconv.Atoi(conns); err == nil && parsed > 0 {
			maxIdleConns = parsed
		}
	}

	return &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "127.0.0.1"),
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			IdleTimeout:     idleTimeout,
			ShutdownTimeout: shutdownTimeout,
			MaxRequestBody:  maxRequestBody,
		},
		Database: DatabaseConfig{
			Path:         dbPath,
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", ""),
		},
	}, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration парсит строку длительности (например, "15s", "1m").
func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}
