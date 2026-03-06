// Package database предоставляет функции для работы с базой данных SQLite.
package database

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// DB обёртка над sql.DB для управления базой данных.
type DB struct {
	*sql.DB
	logger *zap.Logger
}

// New создаёт новое соединение с базой данных.
func New(dbPath string, maxOpenConns, maxIdleConns int, logger *zap.Logger) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	// Включаем WAL режим для лучшей производительности
	_, err = db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA synchronous=NORMAL;
		PRAGMA cache_size=10000;
	`)
	if err != nil {
		logger.Warn("failed to set PRAGMA options", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established", zap.String("path", dbPath))
	return &DB{db, logger}, nil
}

// Init создаёт таблицы и индексы при старте приложения.
func (db *DB) Init() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_name ON users(name)`,
		`CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	db.logger.Info("database initialized successfully")
	return nil
}

// Close закрывает соединение с базой данных.
func (db *DB) Close() error {
	db.logger.Info("closing database connection")
	
	// Чекпоинт WAL перед закрытием
	_, _ = db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	
	return db.DB.Close()
}

// IsUniqueConstraintError проверяет, является ли ошибка нарушением уникальности.
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	sqliteErr, ok := err.(sqlite3.Error)
	return ok && sqliteErr.Code == sqlite3.ErrConstraint
}
