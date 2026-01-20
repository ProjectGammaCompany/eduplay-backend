package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PostEventBlock(ctx context.Context, req *eventModel.PostEventBlockIn) (string, error) {
	const op = "event.UseCase.PostEventBlock"

	s.log.With(slog.String("op", op)).Info("attempting to post event block")

	blockDto := eventModel.PostEventBlockToDto(req)

	ret, err := s.eventClient.PostEventBlock(ctx, blockDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post event block", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("event block posted", slog.Any("event", ret))

	return ret.Message, nil
}
