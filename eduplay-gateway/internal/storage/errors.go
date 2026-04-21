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
	ErrUserIsNotPlayer      = errors.New("user is not player")
	ErrCodeExpired          = errors.New("code expired")

	ErrInfoSegmentAnswerIncorrect    = errors.New("info segment answer is incorrect")
	ErrSingleChoiceAnswerIncorrect   = errors.New("single choice answer is incorrect")
	ErrMultipleChoiceAnswerIncorrect = errors.New("multiple choice answer is incorrect")
	ErrTextAnswerIncorrect           = errors.New("text answer is incorrect")

	ErrNoRows           = errors.New("no rows were found")
	ErrNotFound         = errors.New("not found")
	ErrInvalidOperation = errors.New("incorrect user action")

	ErrJoinCodeRetryFailed = errors.New("failed to generate and save join code")
	ErrEventIsPrivate      = errors.New("event is private")
	ErrEventHasNoGroups    = errors.New("event has no groups")
)
