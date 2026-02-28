package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutEventBlock(ctx context.Context, req *eventModel.PutEventBlockIn) (string, error) {
	const op = "event.UseCase.PutEventBlock"

	s.log.With(slog.String("op", op)).Info("attempting to update event block")

	blockDto := eventModel.PutEventBlockToDto(req)

	ret, err := s.eventClient.PutEventBlock(ctx, blockDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to update event block", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("event block updated", slog.Any("event", ret))

	return ret.Message, nil
}
