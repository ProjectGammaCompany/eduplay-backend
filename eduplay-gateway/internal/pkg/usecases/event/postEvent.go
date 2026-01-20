package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PostEvent(ctx context.Context, req *eventModel.PostEventIn) (string, error) {
	const op = "event.UseCase.PostEvent"

	s.log.With(slog.String("op", op)).Info("attempting to post event")

	eventDto, err := eventModel.PostEventInToDto(req)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to convert event to dto", err.Error())
		return "", err
	}

	ret, err := s.eventClient.PostEvent(ctx, eventDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post event", err.Error())
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("event posted", slog.Any("event", ret))

	return ret.Message, nil
}
