// Package models содержит ошибки валидации.
package models

import "errors"

var (
	// ErrEmptyName возвращается, когда имя пользователя пустое.
	ErrEmptyName = errors.New("имя пользователя обязательно и должно содержать хотя бы один символ")
	
	// ErrNameTooLong возвращается, когда имя превышает 100 символов.
	ErrNameTooLong = errors.New("имя не должно превышать 100 символов")
	
	// ErrUserNotFound возвращается, когда пользователь не найден.
	ErrUserNotFound = errors.New("пользователь не найден")
	
	// ErrUserExists возвращается при попытке создать существующего пользователя.
	ErrUserExists = errors.New("пользователь с таким именем уже существует")
)
