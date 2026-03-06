package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestSuccess тестирует функцию Success.
func TestSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     int
		data       interface{}
		message    string
		wantStatus int
	}{
		{
			name:       "with message",
			status:     http.StatusOK,
			data:       map[string]string{"key": "value"},
			message:    "Success",
			wantStatus: http.StatusOK,
		},
		{
			name:       "without message",
			status:     http.StatusCreated,
			data:       map[string]int64{"id": 1},
			message:    "",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "nil data",
			status:     http.StatusNoContent,
			data:       nil,
			message:    "",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			Success(w, tt.status, tt.data, tt.message)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			// Проверяем Content-Type
			if ct := w.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", ct)
			}

			// Проверяем структуру ответа
			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			// Проверяем наличие timestamp
			if _, ok := resp["timestamp"]; !ok {
				t.Error("expected timestamp in response")
			}

			// Проверяем message если есть
			if tt.message != "" {
				if msg, ok := resp["message"].(string); !ok || msg != tt.message {
					t.Errorf("expected message %q, got %v", tt.message, resp["message"])
				}
			}
		})
	}
}

// TestError тестирует функцию Error.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     int
		message    string
		wantStatus int
	}{
		{
			name:       "bad request",
			status:     http.StatusBadRequest,
			message:    "Invalid input",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			status:     http.StatusNotFound,
			message:    "Resource not found",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "internal error",
			status:     http.StatusInternalServerError,
			message:    "Internal server error",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			Error(w, tt.status, tt.message)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			// Проверяем структуру ответа
			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			// Проверяем наличие error
			if errMsg, ok := resp["error"].(string); !ok || errMsg != tt.message {
				t.Errorf("expected error %q, got %v", tt.message, resp["error"])
			}

			// Проверяем наличие timestamp
			if _, ok := resp["timestamp"]; !ok {
				t.Error("expected timestamp in response")
			}
		})
	}
}

// TestResponseTimestampFormat проверяет формат timestamp.
func TestResponseTimestampFormat(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	Success(w, http.StatusOK, map[string]string{"test": "value"}, "")

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	timestamp, ok := resp["timestamp"].(string)
	if !ok {
		t.Fatal("timestamp is not a string")
	}

	// Проверяем что timestamp в формате RFC3339
	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Errorf("timestamp is not in RFC3339 format: %v", err)
	}
}

// TestJSONEncoding проверяет корректность JSON кодирования.
func TestJSONEncoding(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"id":   int64(1),
		"name": "test",
		"nested": map[string]string{
			"key": "value",
		},
	}

	Success(w, http.StatusOK, data, "")

	// Проверяем что response валидный JSON
	if !json.Valid(w.Body.Bytes()) {
		t.Error("response is not valid JSON")
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	respData, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data is not an object")
	}

	if respData["id"].(float64) != 1 {
		t.Error("id mismatch")
	}
	if respData["name"].(string) != "test" {
		t.Error("name mismatch")
	}
}
