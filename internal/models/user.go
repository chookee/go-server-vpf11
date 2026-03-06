// Package models определяет структуры данных приложения.
package models

import "time"

// User представляет пользователя в системе.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest представляет запрос на создание пользователя.
type CreateUserRequest struct {
	Name string `json:"name"`
}

// UpdateUserRequest представляет запрос на обновление пользователя.
type UpdateUserRequest struct {
	Name string `json:"name"`
}

// Validate проверяает корректность данных пользователя.
func (r *CreateUserRequest) Validate() error {
	if r.Name == "" {
		return ErrEmptyName
	}
	if len(r.Name) > 100 {
		return ErrNameTooLong
	}
	return nil
}

// Validate проверяает корректность данных для обновления.
func (r *UpdateUserRequest) Validate() error {
	if r.Name == "" {
		return ErrEmptyName
	}
	if len(r.Name) > 100 {
		return ErrNameTooLong
	}
	return nil
}

// UsersResponse представляет ответ со списком пользователей.
type UsersResponse struct {
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pagination"`
}

// Pagination содержит информацию о пагинации.
type Pagination struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	TotalPages int `json:"pages"`
}

// StatsResponse представляет статистику приложения.
type StatsResponse struct {
	ActiveConnections int `json:"active_connections"`
	TotalUsers        int `json:"total_users"`
}
