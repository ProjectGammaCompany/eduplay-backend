package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) DeleteBlockCondition(ctx context.Context, req *eventModel.Id) (string, error) {
	const op = "event.UseCase.DeleteBlockCondition"

	s.log.With(slog.String("op", op)).Info("attempting to delete block condition by id")

	ret, err := s.eventClient.DeleteBlockCondition(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to delete block condition by id", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("deleted block condition", slog.Any("condition", req.Id))

	return ret.Message, nil
}
