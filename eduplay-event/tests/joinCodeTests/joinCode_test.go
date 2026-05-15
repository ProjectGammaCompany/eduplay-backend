package tests

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	dto "eduplay-event/internal/generated"
	event "eduplay-event/internal/pkg/usecase/events"
	"eduplay-event/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errs "eduplay-event/internal/storage"
)

func TestEventUserRating_PostJoinCodeUnsuccessful(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)

	mockStorage.On("InsertJoinCode", ctx, "event1", mock.Anything).
		Return(nil, errs.ErrJoinCodeNotUnique)

	mockStorage.On("InsertJoinCode", ctx, "event1", mock.Anything).
		Return(nil, errs.ErrJoinCodeNotUnique)

	mockStorage.On("InsertJoinCode", ctx, "event1", mock.Anything).
		Return(nil, errs.ErrJoinCodeNotUnique)

	mockStorage.On("InsertJoinCode", ctx, "event1", mock.Anything).
		Return(nil, errs.ErrJoinCodeNotUnique)

	mockStorage.On("InsertJoinCode", ctx, "event1", mock.Anything).
		Return(nil, errs.ErrJoinCodeNotUnique)

	mockStorage.On("InsertJoinCode", ctx, "event1", mock.Anything).
		Return(nil, errs.ErrJoinCodeNotUnique)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockStorage, "secret")

	res, err := uc.PostJoinCode(ctx, &dto.Id{Id: "event1"})

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.ErrorIs(t, err, errs.ErrJoinCodeRetryFailed)

	mockStorage.AssertExpectations(t)
}

func TestEventUserRating_Success(t *testing.T) {
	ctx := context.Background()

	returnTime := time.Now().UTC().Add(4 * time.Hour)

	mockStorage := new(mocks.Storage)

	mockStorage.On("InsertJoinCode", ctx, "event1", mock.Anything).
		Return(&returnTime, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockStorage, "secret")

	res, err := uc.PostJoinCode(ctx, &dto.Id{Id: "event1"})

	assert.NoError(t, err)
	assert.NotNil(t, res)

	mockStorage.AssertExpectations(t)
}
