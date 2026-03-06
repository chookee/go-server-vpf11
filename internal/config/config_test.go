// Package config предоставляет конфигурацию приложения.
package config

import (
	"os"
	"testing"
	"time"
)

// TestLoad_Defaults тестирует загрузку конфигурации со значениями по умолчанию.
func TestLoad_Defaults(t *testing.T) {
	t.Parallel()

	// Сохраняем текущие переменные окружения
	originalEnv := saveEnv()
	defer restoreEnv(originalEnv)

	// Очищаем переменные окружения
	clearEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем значения по умолчанию
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("expected port 8080, got %s", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 15*time.Second {
		t.Errorf("expected read timeout 15s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 15*time.Second {
		t.Errorf("expected write timeout 15s, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.Server.IdleTimeout != 60*time.Second {
		t.Errorf("expected idle timeout 60s, got %v", cfg.Server.IdleTimeout)
	}
	if cfg.Server.ShutdownTimeout != 30*time.Second {
		t.Errorf("expected shutdown timeout 30s, got %v", cfg.Server.ShutdownTimeout)
	}
	if cfg.Server.MaxRequestBody != 1024*1024 {
		t.Errorf("expected max request body 1MB, got %d", cfg.Server.MaxRequestBody)
	}
	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("expected max open connections 25, got %d", cfg.Database.MaxOpenConns)
	}
	if cfg.Database.MaxIdleConns != 5 {
		t.Errorf("expected max idle connections 5, got %d", cfg.Database.MaxIdleConns)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected log level info, got %s", cfg.Log.Level)
	}
}

// TestLoad_FromEnv тестирует загрузку конфигурации из переменных окружения.
func TestLoad_FromEnv(t *testing.T) {
	t.Parallel()

	originalEnv := saveEnv()
	defer restoreEnv(originalEnv)

	// Устанавливаем тестовые переменные
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("READ_TIMEOUT", "30s")
	os.Setenv("WRITE_TIMEOUT", "30s")
	os.Setenv("IDLE_TIMEOUT", "120s")
	os.Setenv("SHUTDOWN_TIMEOUT", "60s")
	os.Setenv("MAX_REQUEST_BODY", "2097152")
	os.Setenv("DB_MAX_OPEN_CONNS", "50")
	os.Setenv("DB_MAX_IDLE_CONNS", "10")
	os.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected host 0.0.0.0, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("expected read timeout 30s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.MaxRequestBody != 2097152 {
		t.Errorf("expected max request body 2MB, got %d", cfg.Server.MaxRequestBody)
	}
	if cfg.Database.MaxOpenConns != 50 {
		t.Errorf("expected max open connections 50, got %d", cfg.Database.MaxOpenConns)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log level debug, got %s", cfg.Log.Level)
	}
}

// TestParseDuration тестирует парсинг длительности.
func TestParseDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		expected  time.Duration
		expectErr bool
	}{
		{"seconds", "15s", 15 * time.Second, false},
		{"minutes", "1m", time.Minute, false},
		{"combined", "1m30s", 90 * time.Second, false},
		{"milliseconds", "500ms", 500 * time.Millisecond, false},
		{"hours", "1h", time.Hour, false},
		{"invalid", "invalid", 0, true},
		{"empty", "", 0, true},
		{"negative", "-10s", -10 * time.Second, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d, err := parseDuration(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("error = %v, wantErr %v", err, tt.expectErr)
			}
			if d != tt.expected {
				t.Errorf("got %v, want %v", d, tt.expected)
			}
		})
	}
}

// TestGetEnv тестирует получение переменных окружения.
func TestGetEnv(t *testing.T) {
	t.Parallel()

	originalEnv := saveEnv()
	defer restoreEnv(originalEnv)

	// Тест с установленной переменной
	os.Setenv("TEST_VAR", "test_value")
	if got := getEnv("TEST_VAR", "default"); got != "test_value" {
		t.Errorf("expected test_value, got %s", got)
	}

	// Тест с переменной по умолчанию
	os.Unsetenv("TEST_VAR")
	if got := getEnv("TEST_VAR", "default"); got != "default" {
		t.Errorf("expected default, got %s", got)
	}
}

// saveEnv сохраняет текущие переменные окружения.
func saveEnv() map[string]string {
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT",
		"READ_TIMEOUT", "WRITE_TIMEOUT", "IDLE_TIMEOUT", "SHUTDOWN_TIMEOUT",
		"MAX_REQUEST_BODY",
		"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS",
		"LOG_LEVEL", "LOG_FILE",
		"DB_PATH",
	}

	result := make(map[string]string)
	for _, v := range envVars {
		result[v] = os.Getenv(v)
	}
	return result
}

// restoreEnv восстанавливает переменные окружения.
func restoreEnv(env map[string]string) {
	for k, v := range env {
		if v == "" {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}
}

// clearEnv очищает переменные окружения.
func clearEnv() {
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT",
		"READ_TIMEOUT", "WRITE_TIMEOUT", "IDLE_TIMEOUT", "SHUTDOWN_TIMEOUT",
		"MAX_REQUEST_BODY",
		"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS",
		"LOG_LEVEL", "LOG_FILE",
		"DB_PATH",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}
}
