package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetGroups(ctx context.Context, req *eventModel.Id) (*eventModel.GroupsShort, error) {
	const op = "event.UseCase.GetGroups"

	s.log.With(slog.String("op", op)).Info("attempting to get event groups")

	eventId := &eventDto.Id{Id: req.Id}

	groups, err := s.eventClient.GetGroups(ctx, eventId)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get groups", slog.String("error", err.Error()))
		return nil, err
	}

	return eventModel.GroupsShortFromDto(groups), nil
}
