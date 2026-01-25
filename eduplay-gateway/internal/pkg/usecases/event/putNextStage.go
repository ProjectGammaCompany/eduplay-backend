package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutNextStage(ctx context.Context, in *eventModel.EventBlockTaskUserIds) (string, error) {
	const op = "event.UseCase.PutNextStage"

	s.log.With(slog.String("op", op)).Info("attempting to put next stage")

	ret, err := s.eventClient.PutNextStage(ctx, eventModel.EventBlockTaskUserIdsToDto(in))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to put next stage", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("event put next stage", slog.Any("event", ret))

	return ret.Message, nil
}
