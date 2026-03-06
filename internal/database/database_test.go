package database

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

// setupTestDB создаёт тестовую базу данных в памяти.
func setupTestDB(t *testing.T) (*DB, func()) {
	t.Helper()

	logger, _ := zap.NewDevelopment()

	// Создаём временный файл для БД
	tmpFile, err := os.CreateTemp("", "test_database_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	db, err := New(dbPath, 5, 2, logger)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	if err := db.Init(); err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	// Функция очистки
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, cleanup
}

// TestNew тестирует создание соединения с БД.
func TestNew(t *testing.T) {
	t.Parallel()

	logger, _ := zap.NewDevelopment()

	// Тест с временным файлом
	tmpFile, err := os.CreateTemp("", "test_database_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	db, err := New(dbPath, 5, 2, logger)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer db.Close()

	if db.DB == nil {
		t.Error("expected db.DB to be initialized")
	}
}

// TestNew_InvalidPath тестирует ошибку при неверном пути.
func TestNew_InvalidPath(t *testing.T) {
	t.Parallel()

	logger, _ := zap.NewDevelopment()

	// Путь к несуществующей директории
	_, err := New("/nonexistent/path/to/db.sqlite", 5, 2, logger)
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

// TestInit тестирует инициализацию схемы БД.
func TestInit(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Проверяем что таблица создана
	var tableName string
	err := db.QueryRow(`
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name='users'
	`).Scan(&tableName)

	if err != nil {
		t.Errorf("expected users table to exist: %v", err)
	}

	// Проверяем индексы
	var indexCount int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master 
		WHERE type='index' AND tbl_name='users'
	`).Scan(&indexCount)

	if err != nil {
		t.Errorf("failed to count indexes: %v", err)
	}
	if indexCount < 2 {
		t.Errorf("expected at least 2 indexes, got %d", indexCount)
	}
}

// TestDB_CreateUser тестирует создание пользователя.
func TestDB_CreateUser(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	result, err := db.Exec("INSERT INTO users(name) VALUES (?)", "Alice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("expected no error getting last insert id: %v", err)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
}

// TestDB_CreateUser_Duplicate тестирует уникальность имени.
func TestDB_CreateUser_Duplicate(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Создаём первого пользователя
	_, err := db.Exec("INSERT INTO users(name) VALUES (?)", "Alice")
	if err != nil {
		t.Fatalf("failed to insert first user: %v", err)
	}

	// Пытаемся создать дубликат
	_, err = db.Exec("INSERT INTO users(name) VALUES (?)", "Alice")
	if err == nil {
		t.Error("expected error for duplicate name, got nil")
	}

	// Проверяем что это ошибка уникальности
	if !IsUniqueConstraintError(err) {
		t.Error("expected unique constraint error")
	}
}

// TestDB_GetUser тестирует получение пользователя.
func TestDB_GetUser(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Создаём пользователя
	_, err := db.Exec("INSERT INTO users(name) VALUES (?)", "Bob")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Получаем пользователя
	var id int64
	var name string
	err = db.QueryRow("SELECT id, name FROM users WHERE name = ?", "Bob").Scan(&id, &name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if name != "Bob" {
		t.Errorf("expected name Bob, got %s", name)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
}

// TestDB_UpdateUser тестирует обновление пользователя.
func TestDB_UpdateUser(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Создаём пользователя
	_, err := db.Exec("INSERT INTO users(name) VALUES (?)", "Charlie")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Обновляем
	result, err := db.Exec("UPDATE users SET name = ? WHERE id = ?", "Charles", 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 row affected, got %d", rows)
	}

	// Проверяем обновление
	var name string
	err = db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)
	if err != nil {
		t.Fatalf("failed to get updated user: %v", err)
	}
	if name != "Charles" {
		t.Errorf("expected name Charles, got %s", name)
	}
}

// TestDB_DeleteUser тестирует удаление пользователя.
func TestDB_DeleteUser(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Создаём пользователя
	_, err := db.Exec("INSERT INTO users(name) VALUES (?)", "Delete Me")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Удаляем
	result, err := db.Exec("DELETE FROM users WHERE id = ?", 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 row affected, got %d", rows)
	}

	// Проверяем что удалён
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count users: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 users, got %d", count)
	}
}

// TestDB_CountUsers тестирует подсчёт пользователей.
func TestDB_CountUsers(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Создаём нескольких пользователей
	names := []string{"User1", "User2", "User3"}
	for _, name := range names {
		_, err := db.Exec("INSERT INTO users(name) VALUES (?)", name)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}
	}

	// Считаем
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 users, got %d", count)
	}
}

// TestIsUniqueConstraintError тестирует функцию проверки ошибок уникальности.
func TestIsUniqueConstraintError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "other error",
			err:      os.ErrNotExist,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := IsUniqueConstraintError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestDB_Close тестирует закрытие соединения.
func TestDB_Close(t *testing.T) {
	t.Parallel()

	db, _ := setupTestDB(t)

	err := db.Close()
	if err != nil {
		t.Errorf("expected no error on close, got %v", err)
	}
}

// BenchmarkDB_InsertUser бенчмарк вставки пользователя.
func BenchmarkDB_InsertUser(b *testing.B) {
	logger, _ := zap.NewDevelopment()

	tmpFile, _ := os.CreateTemp("", "bench_database_*.db")
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	db, _ := New(dbPath, 5, 2, logger)
	defer db.Close()
	db.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Exec("INSERT INTO users(name) VALUES (?)", "BenchmarkUser")
	}
}
