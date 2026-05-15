package tests

import (
	"context"
	"io"
	"log/slog"
	"testing"

	dto "eduplay-event/internal/generated"
	event "eduplay-event/internal/pkg/usecase/events"
	"eduplay-event/tests/mocks"

	"github.com/stretchr/testify/assert"

	errs "eduplay-event/internal/storage"
)

func TestEventUserRating_NoRating(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)

	mockStorage.On("GetEventUserRating", ctx, "user1", "event1").
		Return(int64(0), errs.ErrNotFound)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockStorage, "secret")

	res, err := uc.GetEventUserRating(ctx, &dto.UserEventIds{
		UserId:  "user1",
		EventId: "event1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "-1", res.Message)

	mockStorage.AssertExpectations(t)
}

func TestEventUserRating_Success(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)

	mockStorage.On("GetEventUserRating", ctx, "user1", "event1").
		Return(int64(4), nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockStorage, "secret")

	res, err := uc.GetEventUserRating(ctx, &dto.UserEventIds{
		UserId:  "user1",
		EventId: "event1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "4", res.Message)

	mockStorage.AssertExpectations(t)
}
