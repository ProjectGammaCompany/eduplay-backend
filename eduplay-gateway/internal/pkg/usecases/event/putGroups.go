package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutGroups(ctx context.Context, req *eventModel.PutGroupsIn) (string, error) {
	const op = "event.UseCase.PutGroups"

	s.log.With(slog.String("op", op)).Info("attempting to update event groups")

	message, err := s.eventClient.PutGroups(ctx, eventModel.PutGroupsInToDto(req))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to update groups", slog.String("error", err.Error()))
		return "", err
	}

	return message.Message, nil
}
