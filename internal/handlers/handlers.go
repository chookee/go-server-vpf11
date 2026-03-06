// Package handlers содержит HTTP обработчики для API endpoints.
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/zerocode/users-api/internal/database"
	"github.com/zerocode/users-api/internal/models"
	"github.com/zerocode/users-api/internal/utils"
	"go.uber.org/zap"
)

// Handler содержит зависимости для обработчиков.
type Handler struct {
	db     *sql.DB
	logger *zap.Logger
}

// New создаёт новый экземпляр Handler.
func New(db *sql.DB, logger *zap.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
	}
}

// HealthCheck возвращает статус здоровья приложения.
// @route GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "healthy"}
	utils.Success(w, http.StatusOK, data, "")
}

// GetUsers возвращает список всех пользователей с пагинацией.
// @route GET /users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Парсим параметры пагинации с валидацией
	page := parseIntParam(r.URL.Query().Get("page"), 1)
	if page < 1 {
		page = 1
	}
	perPage := parseIntParam(r.URL.Query().Get("per_page"), 10)
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	
	// Вычисляем offset безопасно (без переполнения)
	offset := (page - 1) * perPage

	// Получаем пользователей
	rows, err := h.db.QueryContext(r.Context(),
		"SELECT id, name, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?",
		perPage, offset,
	)
	if err != nil {
		h.logger.Error("failed to query users", zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.CreatedAt); err != nil {
			h.logger.Error("failed to scan user", zap.Error(err))
			utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		h.logger.Error("rows iteration error", zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	// Получаем общее количество
	var total int
	err = h.db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		h.logger.Error("failed to count users", zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	totalPages := (total + perPage - 1) / perPage
	if totalPages == 0 {
		totalPages = 1
	}

	response := models.UsersResponse{
		Users: users,
		Pagination: models.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	utils.Success(w, http.StatusOK, response, "")
}

// AddUser создаёт нового пользователя.
// @route POST /add-user
func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := parseJSON(w, r, &req); err != nil {
		return
	}

	// Нормализуем имя (trim)
	req.Name = strings.TrimSpace(req.Name)

	if err := req.Validate(); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.db.ExecContext(r.Context(),
		"INSERT INTO users(name) VALUES (?)", req.Name)
	if err != nil {
		if database.IsUniqueConstraintError(err) {
			h.logger.Info("duplicate user", zap.String("name", req.Name))
			utils.Error(w, http.StatusConflict, models.ErrUserExists.Error())
			return
		}
		h.logger.Error("failed to insert user", zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	id, _ := result.LastInsertId()

	h.logger.Info("user created", zap.Int64("id", id), zap.String("name", req.Name))

	user := models.User{
		ID:   id,
		Name: req.Name,
	}

	utils.Success(w, http.StatusCreated, user, "Пользователь успешно добавлен")
}

// GetUser возвращает пользователя по ID.
// @route GET /user/{id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var u models.User
	err = h.db.QueryRowContext(r.Context(),
		"SELECT id, name, created_at FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Name, &u.CreatedAt)

	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusNotFound, models.ErrUserNotFound.Error())
		return
	}
	if err != nil {
		h.logger.Error("failed to get user", zap.Int64("id", id), zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	utils.Success(w, http.StatusOK, u, "")
}

// UpdateUser обновляет данные пользователя.
// @route PUT /user/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var req models.UpdateUserRequest
	if err := parseJSON(w, r, &req); err != nil {
		return
	}

	// Нормализуем имя
	req.Name = strings.TrimSpace(req.Name)

	if err := req.Validate(); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.db.ExecContext(r.Context(),
		"UPDATE users SET name = ? WHERE id = ?", req.Name, id)
	if err != nil {
		if database.IsUniqueConstraintError(err) {
			h.logger.Info("duplicate user on update", zap.String("name", req.Name))
			utils.Error(w, http.StatusConflict, models.ErrUserExists.Error())
			return
		}
		h.logger.Error("failed to update user", zap.Int64("id", id), zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		utils.Error(w, http.StatusNotFound, models.ErrUserNotFound.Error())
		return
	}

	h.logger.Info("user updated", zap.Int64("id", id), zap.String("name", req.Name))

	user := models.User{
		ID:   id,
		Name: req.Name,
	}

	utils.Success(w, http.StatusOK, user, "Данные пользователя обновлены")
}

// DeleteUser удаляет пользователя.
// @route DELETE /user/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	result, err := h.db.ExecContext(r.Context(),
		"DELETE FROM users WHERE id = ?", id)
	if err != nil {
		h.logger.Error("failed to delete user", zap.Int64("id", id), zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		utils.Error(w, http.StatusNotFound, models.ErrUserNotFound.Error())
		return
	}

	h.logger.Info("user deleted", zap.Int64("id", id))

	utils.Success(w, http.StatusOK, map[string]int64{"id": id}, "Пользователь успешно удалён")
}

// GetStats возвращает статистику приложения.
// @route GET /stats/active
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	var total int
	err := h.db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		h.logger.Error("failed to get stats", zap.Error(err))
		utils.Error(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	stats := models.StatsResponse{
		ActiveConnections: 0,
		TotalUsers:        total,
	}

	utils.Success(w, http.StatusOK, stats, "")
}

// =============================================================================
// Вспомогательные функции
// =============================================================================

// parseJSON парсит JSON из тела запроса.
func parseJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.Body == nil {
		utils.Error(w, http.StatusBadRequest, "Тело запроса обязательно")
		return errors.New("empty body")
	}
	defer r.Body.Close()
	
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		utils.Error(w, http.StatusBadRequest, "Некорректный JSON")
		return err
	}
	return nil
}

// parseIntParam парсит целочисленный параметр со значением по умолчанию.
func parseIntParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}
