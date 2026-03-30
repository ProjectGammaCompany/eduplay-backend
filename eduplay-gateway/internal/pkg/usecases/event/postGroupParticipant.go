package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	errs "eduplay-gateway/internal/storage"
	"fmt"
	"log/slog"
)

func (s *UseCase) PostGroupParticipant(ctx context.Context, req *eventModel.ParticipationPasswords) (string, error) {
	const op = "event.UseCase.PostGroupParticipant"

	s.log.With(slog.String("op", op)).Info("attempting to post group participant")

	ret, err := s.eventClient.GetEvent(ctx, &eventDto.Id{Id: req.EventId})

	if err != nil {
		if err == errs.ErrNotFound {
			return "", errs.ErrNotFound
		}
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return "", err
	}

	if ret.Private {
		return "", errs.ErrEventIsPrivate
	}

	if !ret.GroupEvent {
		return "", errs.ErrEventHasNoGroups
	}

	role, err := s.GetRole(ctx, req.UserId, req.EventId)
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
	if ret.GroupEvent {
		eventGroups, err := s.eventClient.GetGroups(ctx, &eventDto.Id{Id: req.EventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get groups", slog.String("error", err.Error()))
			return "", err
		}
		correctLoginInfo := false
		for _, eventGroup := range eventGroups.Groups {
			fmt.Println(eventGroup.Login, eventGroup.Password)
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

	s.log.With(slog.String("op", op)).Info("get event by join code", slog.Any("event", ret))

	message, err := s.eventClient.PostParticipant(ctx, &eventDto.PostParticipantIn{UserId: req.UserId, EventId: req.EventId, GroupId: groupId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post participant", slog.String("error", err.Error()))
		return "", err
	}

	return message.Message, nil
}
