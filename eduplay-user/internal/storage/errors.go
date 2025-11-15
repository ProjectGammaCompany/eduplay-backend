package storage

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidRefresh    = errors.New("invalid refresh token")
	ErrIsActive          = errors.New("user is already active")
)
