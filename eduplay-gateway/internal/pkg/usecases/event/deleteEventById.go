package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) DeleteEventById(ctx context.Context, req *eventModel.Id) (string, error) {
	const op = "event.UseCase.DeleteEventById"

	s.log.With(slog.String("op", op)).Info("attempting to delete event by id")

	ret, err := s.eventClient.DeleteEventById(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to delete event by id", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("deleted event", slog.Any("event", req.Id))

	return ret.Message, nil
}
