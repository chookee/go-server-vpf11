// Package main - точка входа приложения.
package main

import (
	"testing"

	"github.com/zerocode/users-api/internal/config"
	"go.uber.org/zap"
)

// TestInitLogger тестирует инициализацию логгера.
func TestInitLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     config.LogConfig
		wantErr bool
	}{
		{
			name: "info level",
			cfg: config.LogConfig{
				Level: "info",
				File:  "",
			},
			wantErr: false,
		},
		{
			name: "debug level",
			cfg: config.LogConfig{
				Level: "debug",
				File:  "",
			},
			wantErr: false,
		},
		{
			name: "production level",
			cfg: config.LogConfig{
				Level: "production",
				File:  "",
			},
			wantErr: false,
		},
		{
			name: "default level",
			cfg: config.LogConfig{
				Level: "",
				File:  "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger, err := initLogger(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("initLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if logger == nil && !tt.wantErr {
				t.Error("initLogger() logger = nil, want not nil")
			}
			if logger != nil {
				logger.Sync()
			}
		})
	}
}

// TestInitLogger_WithFile тестирует инициализацию логгера с файлом.
func TestInitLogger_WithFile(t *testing.T) {
	t.Parallel()

	cfg := config.LogConfig{
		Level: "debug",
		File:  "", // Пустой файл = stdout
	}

	logger, err := initLogger(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("expected logger, got nil")
	}

	// Проверяем что логгер работает
	logger.Info("test message")
}

// TestInitLogger_Levels тестирует разные уровни логирования.
func TestInitLogger_Levels(t *testing.T) {
	t.Parallel()

	levels := []string{"debug", "info", "production", ""}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			t.Parallel()

			cfg := config.LogConfig{
				Level: level,
				File:  "",
			}

			logger, err := initLogger(cfg)
			if err != nil {
				t.Fatalf("level %s: expected no error, got %v", level, err)
			}
			defer logger.Sync()

			if logger == nil {
				t.Fatalf("level %s: expected logger, got nil", level)
			}
		})
	}
}
