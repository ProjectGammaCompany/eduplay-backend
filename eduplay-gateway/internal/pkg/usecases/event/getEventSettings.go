package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetEventSettings(ctx context.Context, req *eventModel.Id) (*eventModel.GetEventSettings, error) {
	const op = "event.UseCase.GetEventSettings"

	s.log.With(slog.String("op", op)).Info("attempting to get event settings")

	eventId := &eventDto.Id{Id: req.Id}

	eventInfo, err := s.eventClient.GetEvent(ctx, eventId)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", err.Error())
		return nil, err
	}

	groups, err := s.eventClient.GetGroups(ctx, eventId)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get groups", err.Error())
		return nil, err
	}

	collaborators, err := s.eventClient.GetCollaborators(ctx, eventId)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get collaborators", err.Error())
		return nil, err
	}

	ret := eventModel.GetEventSettingsFromDto(eventInfo, groups, collaborators)

	s.log.With(slog.String("op", op)).Info("got event", slog.Any("event", ret))

	return ret, nil
}
