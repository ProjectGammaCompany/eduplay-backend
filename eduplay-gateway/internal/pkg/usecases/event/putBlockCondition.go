package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutBlockCondition(ctx context.Context, req *eventModel.Condition) (string, error) {
	const op = "event.UseCase.PutBlockCondition"

	s.log.With(slog.String("op", op)).Info("attempting to update block condition")

	taskDto := eventModel.ConditionToDto(req)

	ret, err := s.eventClient.PutBlockCondition(ctx, taskDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to update block condition", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("block condition updated", slog.Any("event", ret))

	return ret.Message, nil
}
