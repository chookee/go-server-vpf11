// Package utils предоставляет вспомогательные функции для API.
package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Response представляет стандартный ответ API.
type Response struct {
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// JSON устанавливает заголовки и кодирует данные в JSON.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		zap.L().Error("failed to encode response", zap.Error(err))
	}
}

// Success возвращает успешный ответ API.
func Success(w http.ResponseWriter, status int, data interface{}, message string) {
	response := Response{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	if message != "" {
		response.Message = message
	}
	JSON(w, status, response)
}

// Error возвращает ответ с ошибкой API.
func Error(w http.ResponseWriter, status int, message string) {
	response := Response{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Error:     message,
	}
	JSON(w, status, response)
}
