package storage

import (
	"errors"
)

var (
	ErrSegNotFound  = errors.New("segment not found")
	ErrSegNotExists = errors.New("segment not exists")
	ErrSegExists    = errors.New("segment already exists")
	ErrUserNotFound = errors.New("user with such segment not found")
)
