package storage

import (
	"errors"
)

var (
	ErrInvalidRequest  = errors.New("failed to deserialize request")
	ErrValidationError = errors.New("failed to validate request")

	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrPasswordsNotMatch    = errors.New("passwords do not match")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrAccessTokenExpired   = errors.New("access token expired")
	ErrUserIsNotOperator    = errors.New("user is not operator")
	ErrIsActive             = errors.New("user is already active")
)
