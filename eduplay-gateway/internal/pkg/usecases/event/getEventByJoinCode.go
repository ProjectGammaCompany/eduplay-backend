package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	errs "eduplay-gateway/internal/storage"
	"fmt"
	"log/slog"
)

func (s *UseCase) GetEventByJoinCode(ctx context.Context, joinCode string, username string) (bool, error) {
	const op = "event.UseCase.GetEventByJoinCode"

	s.log.With(slog.String("op", op)).Info("attempting to get event by join code")

	ret, err := s.eventClient.GetEventByJoinCode(ctx, &eventDto.Id{Id: joinCode})

	if err != nil {
		if err == errs.ErrNotFound {
			return false, errs.ErrNotFound
		}
		s.log.With(slog.String("op", op)).Error("failed to get event id by join code", slog.String("error", err.Error()))
		return false, err
	}

	event, err := s.eventClient.GetEvent(ctx, &eventDto.Id{Id: ret.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return false, err
	}

	collaborators, err := s.eventClient.GetCollaborators(ctx, &eventDto.Id{Id: ret.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get collaborators", slog.String("error", err.Error()))
		return false, err
	}

	for _, collaborator := range collaborators.Users {
		if collaborator.Email == username {
			return false, errs.ErrUserIsNotPlayer
		}
	}

	s.log.With(slog.String("op", op)).Info("get event by join code", slog.Any("event", ret))

	fmt.Println(event.GroupEvent)

	return event.GroupEvent, err
}
