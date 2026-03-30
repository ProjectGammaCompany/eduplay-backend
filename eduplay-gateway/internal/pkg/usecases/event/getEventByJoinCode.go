package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	errs "eduplay-gateway/internal/storage"
	"log/slog"
)

func (s *UseCase) GetEventByJoinCode(ctx context.Context, joinCode string, userId string) (bool, error) {
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

	role, err := s.GetRole(ctx, userId, event.EventId)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get role", slog.String("error", err.Error()))
		return false, err
	}

	if role == 1 {
		return false, errs.ErrUserIsNotPlayer
	}

	s.log.With(slog.String("op", op)).Info("get event by join code", slog.Any("event", ret))

	return event.GroupEvent, err
}
