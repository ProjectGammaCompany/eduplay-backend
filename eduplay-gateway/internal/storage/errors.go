package storage

import (
	"errors"
)

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrAccessTokenExpired   = errors.New("access token expired")
	ErrUserIsNotOperator    = errors.New("user is not operator")
	ErrIsActive             = errors.New("user is already active")
)
