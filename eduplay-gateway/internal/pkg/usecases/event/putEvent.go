package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutEvent(ctx context.Context, req *eventModel.PutEventIn) (*eventModel.Groups, error) {
	const op = "event.UseCase.PutEvent"

	s.log.With(slog.String("op", op)).Info("attempting to put event")

	userIds := make([]string, len(req.Collaborators))
	for i, collaborator := range req.Collaborators {
		userInfo, err := s.userClient.GetProfileByLogin(ctx, collaborator)
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
		}
		userIds[i] = userInfo.UserId
	}

	req.Collaborators = userIds

	eventDto, err := eventModel.PutEventInToDto(req)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to convert event to dto", slog.String("error", err.Error()))
		return nil, err
	}

	ret, err := s.eventClient.PutEvent(ctx, eventDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to put event", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("event put", slog.Any("event", ret))

	return eventModel.GroupsFromDto(ret), nil
}
