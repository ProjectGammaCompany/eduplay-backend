package tests

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	dto "eduplay-user/internal/generated"
	"eduplay-user/internal/pkg/usecase/user"
	"eduplay-user/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUpUser_Success(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// bcrypt hash неизвестен → mock.Anything
	mockStorage.
		On("SignUpUser", ctx, "test@mail.ru", mock.Anything).
		Return("testUserId", nil)

	// refreshToken генерируется → mock.Anything
	mockStorage.
		On("SaveSession", ctx, "testUserId", mock.Anything).
		Return(nil)

	uc := user.New(logger, mockStorage, nil, nil, "secret")

	session, err := uc.SignUpUser(ctx, &dto.SignUpIn{
		Email:    "test@mail.ru",
		Password: "testpassword",
	})

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.RefreshToken)

	mockStorage.AssertExpectations(t)
}

func TestSignUpUser_SaveSessionError(t *testing.T) {

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockStorage.
		On("SignUpUser", mock.Anything, "test@mail.ru", mock.Anything).
		Return("userId", nil)

	mockStorage.
		On("SaveSession", mock.Anything, "userId", mock.Anything).
		Return(errors.New("db error"))

	uc := user.New(logger, mockStorage, nil, nil, "secret")

	_, err := uc.SignUpUser(context.Background(), &dto.SignUpIn{
		Email:    "test@mail.ru",
		Password: "pass",
	})

	assert.Error(t, err)
}

func TestSignUpUser_UserExists(t *testing.T) {

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockStorage.
		On("SignUpUser", mock.Anything, "test@mail.ru", mock.Anything).
		Return("", errors.New("user exists"))

	uc := user.New(logger, mockStorage, nil, nil, "secret")

	_, err := uc.SignUpUser(context.Background(), &dto.SignUpIn{
		Email:    "test@mail.ru",
		Password: "pass",
	})

	assert.Error(t, err)
}
