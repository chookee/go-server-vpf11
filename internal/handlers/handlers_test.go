package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/zerocode/users-api/internal/database"
	"go.uber.org/zap"
)

// setupTestHandler создаёт тестовый handler с тестовой БД.
func setupTestHandler(t *testing.T) (*Handler, func()) {
	t.Helper()

	logger, _ := zap.NewDevelopment()

	// Создаём временный файл для БД
	tmpFile, err := os.CreateTemp("", "test_handlers_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	db, err := database.New(dbPath, 5, 2, logger)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	if err := db.Init(); err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return New(db.DB, logger), cleanup
}

// createTestUser создаёт тестового пользователя в БД.
func createTestUser(t *testing.T, db *database.DB, name string) int64 {
	t.Helper()

	result, err := db.Exec("INSERT INTO users(name) VALUES (?)", name)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	id, _ := result.LastInsertId()
	return id
}

// =============================================================================
// Health Check Tests
// =============================================================================

// TestHealthCheck тестирует endpoint health check.
func TestHealthCheck(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data in response")
	}
	if data["status"] != "healthy" {
		t.Error("expected status healthy")
	}
}

// =============================================================================
// Add User Tests
// =============================================================================

// TestAddUser_Success тестирует успешное создание пользователя.
func TestAddUser_Success(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	body := `{"name": "Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	if data["name"] != "Alice" {
		t.Errorf("expected name Alice, got %v", data["name"])
	}
}

// TestAddUser_EmptyName тестирует создание пользователя с пустым именем.
func TestAddUser_EmptyName(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	body := `{"name": ""}`
	req := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestAddUser_WhitespaceName тестирует создание пользователя с именем из пробелов.
func TestAddUser_WhitespaceName(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	body := `{"name": "   "}`
	req := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestAddUser_NameTooLong тестирует создание пользователя с длинным именем.
func TestAddUser_NameTooLong(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	longName := bytes.Repeat([]byte("a"), 101)
	body := `{"name": "` + string(longName) + `"}`

	req := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestAddUser_Duplicate тестирует создание дубликата пользователя.
func TestAddUser_Duplicate(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// Создаём первого пользователя
	body := `{"name": "Duplicate"}`
	req := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.AddUser(w, req)

	// Пытаемся создать дубликат
	req2 := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewReader([]byte(body)))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.AddUser(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w2.Code)
	}
}

// TestAddUser_InvalidJSON тестирует отправку невалидного JSON.
func TestAddUser_InvalidJSON(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestAddUser_EmptyBody тестирует отправку пустого тела.
func TestAddUser_EmptyBody(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/add-user", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// =============================================================================
// Get User Tests
// =============================================================================

// TestGetUser_Success тестирует успешное получение пользователя.
func TestGetUser_Success(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp("", "test_handlers_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	logger, _ := zap.NewDevelopment()
	testDB, _ := database.New(dbPath, 5, 2, logger)
	defer testDB.Close()
	testDB.Init()

	userID := createTestUser(t, testDB, "GetUserTest")
	h := New(testDB.DB, logger)

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	if int64(data["id"].(float64)) != userID {
		t.Errorf("expected id %d, got %v", userID, data["id"])
	}
}

// TestGetUser_NotFound тестирует получение несуществующего пользователя.
func TestGetUser_NotFound(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/user/999", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// TestGetUser_InvalidID тестирует получение пользователя с невалидным ID.
func TestGetUser_InvalidID(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/user/abc", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// =============================================================================
// Get Users Tests
// =============================================================================

// TestGetUsers_EmptyList тестирует получение пустого списка.
func TestGetUsers_EmptyList(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	
	// Проверяем что users существует и это массив
	dataData := data["data"].(map[string]interface{})
	users, ok := dataData["users"].([]interface{})
	if !ok {
		// Если nil, это тоже нормально - пустой список
		if dataData["users"] == nil {
			// Пустой список - OK
			return
		}
		t.Fatalf("expected users to be array, got %T", dataData["users"])
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

// TestGetUsers_WithPagination тестирует пагинацию.
func TestGetUsers_WithPagination(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp("", "test_handlers_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	logger, _ := zap.NewDevelopment()
	testDB, _ := database.New(dbPath, 5, 2, logger)
	defer testDB.Close()
	testDB.Init()

	// Создаём 5 пользователей
	for i := 1; i <= 5; i++ {
		createTestUser(t, testDB, "User"+string(rune('0'+i)))
	}

	h := New(testDB.DB, logger)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&per_page=2", nil)
	w := httptest.NewRecorder()

	h.GetUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	pagination := data["pagination"].(map[string]interface{})

	if int(pagination["total"].(float64)) != 5 {
		t.Errorf("expected total 5, got %v", pagination["total"])
	}
	if int(pagination["pages"].(float64)) != 3 {
		t.Errorf("expected 3 pages, got %v", pagination["pages"])
	}
}

// TestGetUsers_InvalidPagination тестирует невалидные параметры пагинации.
func TestGetUsers_InvalidPagination(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/users?page=-1&per_page=999", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// =============================================================================
// Update User Tests
// =============================================================================

// TestUpdateUser_Success тестирует успешное обновление.
func TestUpdateUser_Success(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp("", "test_handlers_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	logger, _ := zap.NewDevelopment()
	testDB, _ := database.New(dbPath, 5, 2, logger)
	defer testDB.Close()
	testDB.Init()

	createTestUser(t, testDB, "Original")
	h := New(testDB.DB, logger)

	body := `{"name": "Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/user/1", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// TestUpdateUser_NotFound тестирует обновление несуществующего пользователя.
func TestUpdateUser_NotFound(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	body := `{"name": "Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/user/999", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// =============================================================================
// Delete User Tests
// =============================================================================

// TestDeleteUser_Success тестирует успешное удаление.
func TestDeleteUser_Success(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp("", "test_handlers_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	logger, _ := zap.NewDevelopment()
	testDB, _ := database.New(dbPath, 5, 2, logger)
	defer testDB.Close()
	testDB.Init()

	createTestUser(t, testDB, "ToDelete")
	h := New(testDB.DB, logger)

	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// TestDeleteUser_NotFound тестирует удаление несуществующего пользователя.
func TestDeleteUser_NotFound(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodDelete, "/user/999", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// =============================================================================
// Stats Tests
// =============================================================================

// TestGetStats тестирует статистику.
func TestGetStats(t *testing.T) {
	t.Parallel()

	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/stats/active", nil)
	w := httptest.NewRecorder()

	handler.GetStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	if data["total_users"] == nil {
		t.Error("expected total_users in response")
	}
}

// =============================================================================
// Helper Tests
// =============================================================================

// TestParseIntParam тестирует вспомогательную функцию.
func TestParseIntParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		defaultVal int
		expected   int
	}{
		{"empty string", "", 10, 10},
		{"valid number", "5", 10, 5},
		{"invalid number", "abc", 10, 10},
		{"negative", "-1", 10, -1},
		{"zero", "0", 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := parseIntParam(tt.input, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
