package storage

import (
	"errors"
)

var (
	ErrJoinCodeNotUnique   = errors.New("join code is not unique")
	ErrJoinCodeRetryFailed = errors.New("failed to generate and save join code")
	ErrJoinCodeExpired     = errors.New("join code expired")

// ErrUserAlreadyExists = errors.New("user already exists")
// ErrUserNotFound      = errors.New("user not found")
// ErrIncorrectPassword = errors.New("incorrect password")
// ErrInvalidRefresh    = errors.New("invalid refresh token")
// ErrIsActive          = errors.New("user is already active")
)
