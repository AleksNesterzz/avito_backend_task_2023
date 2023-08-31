package storage

import (
	"errors"
)

var (
	ErrSegNotFound  = errors.New("segment not found")
	ErrSegExists    = errors.New("segment already exists")
	ErrUserNotFound = errors.New("user not found")
)
