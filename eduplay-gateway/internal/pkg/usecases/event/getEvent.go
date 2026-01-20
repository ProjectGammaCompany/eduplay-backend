package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetEvent(ctx context.Context, req *eventModel.Id) (*eventModel.PostEventIn, error) {
	const op = "event.UseCase.GetEvent"

	s.log.With(slog.String("op", op)).Info("attempting to get event")

	ret, err := s.eventClient.GetEvent(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", err.Error())
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got event", slog.Any("event", ret))

	return eventModel.PostEventInFromDto(ret), nil
}
