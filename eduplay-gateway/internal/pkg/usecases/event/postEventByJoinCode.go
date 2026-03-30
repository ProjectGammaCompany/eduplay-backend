package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	errs "eduplay-gateway/internal/storage"
	"log/slog"
)

func (s *UseCase) PostEventByJoinCode(ctx context.Context, req *eventModel.ParticipationPasswords) (string, error) {
	const op = "event.UseCase.PostEventByJoinCode"

	s.log.With(slog.String("op", op)).Info("attempting to post event by join code")

	ret, err := s.eventClient.GetEventByJoinCode(ctx, &eventDto.Id{Id: req.JoinCode})

	if err != nil {
		if err == errs.ErrNotFound {
			return "", errs.ErrNotFound
		}
		s.log.With(slog.String("op", op)).Error("failed to post event id by join code", slog.String("error", err.Error()))
		return "", err
	}

	event, err := s.eventClient.GetEvent(ctx, &eventDto.Id{Id: ret.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return "", err
	}

	role, err := s.GetRole(ctx, req.UserId, event.EventId)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get role", slog.String("error", err.Error()))
		return "", err
	}

	if role == 1 {
		return "", errs.ErrUserIsNotPlayer
	}

	if role == 0 {
		return "", errs.ErrUserAlreadyExists
	}

	groupId := ""
	if event.GroupEvent {
		eventGroups, err := s.eventClient.GetGroups(ctx, &eventDto.Id{Id: event.EventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get groups", slog.String("error", err.Error()))
			return "", err
		}
		correctLoginInfo := false
		for _, eventGroup := range eventGroups.Groups {
			if eventGroup.Login == req.GroupName && eventGroup.Password == req.GroupPassword {
				correctLoginInfo = true
				groupId = eventGroup.Id
				break
			}
		}

		if !correctLoginInfo {
			return "", errs.ErrIncorrectPassword
		}

	}

	if event.Password != req.Password {
		return "", errs.ErrIncorrectPassword
	}

	s.log.With(slog.String("op", op)).Info("get event by join code", slog.Any("event", ret))

	message, err := s.eventClient.PostParticipant(ctx, &eventDto.PostParticipantIn{UserId: req.UserId, EventId: event.EventId, GroupId: groupId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post participant", slog.String("error", err.Error()))
		return "", err
	}

	return message.Message, nil
}
