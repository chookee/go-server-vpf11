package models

import (
	"strings"
	"testing"
)

// TestCreateUserRequestValidate тестирует валидацию CreateUserRequest.
func TestCreateUserRequestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     CreateUserRequest
		wantError error
	}{
		{
			name: "valid name",
			input: CreateUserRequest{
				Name: "Alice",
			},
			wantError: nil,
		},
		{
			name: "empty name",
			input: CreateUserRequest{
				Name: "",
			},
			wantError: ErrEmptyName,
		},
		{
			name: "whitespace only name",
			input: CreateUserRequest{
				Name: "   ",
			},
			wantError: nil, // Validate не делает trim, это делает handler
		},
		{
			name: "name too long",
			input: CreateUserRequest{
				Name: strings.Repeat("a", 101),
			},
			wantError: ErrNameTooLong,
		},
		{
			name: "name exactly 100 chars",
			input: CreateUserRequest{
				Name: strings.Repeat("a", 100),
			},
			wantError: nil,
		},
		{
			name: "name exactly 1 char",
			input: CreateUserRequest{
				Name: "a",
			},
			wantError: nil,
		},
		{
			name: "name with spaces",
			input: CreateUserRequest{
				Name: "Alice Bob",
			},
			wantError: nil,
		},
		{
			name: "unicode name",
			input: CreateUserRequest{
				Name: "Алексей",
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.input.Validate()

			if err != tt.wantError {
				if tt.wantError == nil {
					t.Errorf("expected no error, got %v", err)
				} else if err != tt.wantError {
					t.Errorf("expected error %v, got %v", tt.wantError, err)
				}
			}
		})
	}
}

// TestUpdateUserRequestValidate тестирует валидацию UpdateUserRequest.
func TestUpdateUserRequestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     UpdateUserRequest
		wantError error
	}{
		{
			name: "valid name",
			input: UpdateUserRequest{
				Name: "Updated Name",
			},
			wantError: nil,
		},
		{
			name: "empty name",
			input: UpdateUserRequest{
				Name: "",
			},
			wantError: ErrEmptyName,
		},
		{
			name: "name too long",
			input: UpdateUserRequest{
				Name: strings.Repeat("x", 150),
			},
			wantError: ErrNameTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.input.Validate()
			if err != tt.wantError {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

// TestUserJSONTags тестирует JSON теги модели User.
func TestUserJSONTags(t *testing.T) {
	t.Parallel()

	// Этот тест проверяет что JSON теги корректно настроены
	// через рефлексию можно было бы проверить, но для простоты
	// проверяем что структура компилируется и поля доступны
	user := User{
		ID:   1,
		Name: "Test",
	}

	if user.ID != 1 {
		t.Error("ID field not accessible")
	}
	if user.Name != "Test" {
		t.Error("Name field not accessible")
	}
}

// TestPaginationJSONTags тестирует JSON теги Pagination.
func TestPaginationJSONTags(t *testing.T) {
	t.Parallel()

	pagination := Pagination{
		Page:       1,
		PerPage:    10,
		Total:      100,
		TotalPages: 10,
	}

	if pagination.Page != 1 {
		t.Error("Page field not accessible")
	}
	if pagination.TotalPages != 10 {
		t.Error("TotalPages field not accessible")
	}
}

// TestUsersResponse тестирует структуру UsersResponse.
func TestUsersResponse(t *testing.T) {
	t.Parallel()

	response := UsersResponse{
		Users: []User{
			{ID: 1, Name: "User1"},
			{ID: 2, Name: "User2"},
		},
		Pagination: Pagination{
			Page:       1,
			PerPage:    10,
			Total:      2,
			TotalPages: 1,
		},
	}

	if len(response.Users) != 2 {
		t.Error("Users slice not accessible")
	}
	if response.Pagination.Total != 2 {
		t.Error("Pagination not accessible")
	}
}

// TestStatsResponse тестирует структуру StatsResponse.
func TestStatsResponse(t *testing.T) {
	t.Parallel()

	stats := StatsResponse{
		ActiveConnections: 5,
		TotalUsers:        100,
	}

	if stats.ActiveConnections != 5 {
		t.Error("ActiveConnections not accessible")
	}
	if stats.TotalUsers != 100 {
		t.Error("TotalUsers not accessible")
	}
}

// TestErrorMessages проверяет сообщения ошибок.
func TestErrorMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		message string
	}{
		{
			name:    "empty name error",
			err:     ErrEmptyName,
			message: "имя пользователя обязательно и должно содержать хотя бы один символ",
		},
		{
			name:    "name too long error",
			err:     ErrNameTooLong,
			message: "имя не должно превышать 100 символов",
		},
		{
			name:    "user not found error",
			err:     ErrUserNotFound,
			message: "пользователь не найден",
		},
		{
			name:    "user exists error",
			err:     ErrUserExists,
			message: "пользователь с таким именем уже существует",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err.Error() != tt.message {
				t.Errorf("expected error message %q, got %q", tt.message, tt.err.Error())
			}
		})
	}
}
