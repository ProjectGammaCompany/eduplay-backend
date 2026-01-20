package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	"log/slog"
)

func (s *UseCase) GetRole(ctx context.Context, userId string, eventId string) (int64, error) {
	const op = "event.UseCase.GetRole"

	s.log.With(slog.String("op", op)).Info("attempting to get role")

	ret, err := s.eventClient.GetRole(ctx, &eventDto.GetRoleIn{UserId: userId, EventId: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get role", err.Error())
		return 0, err
	}

	s.log.With(slog.String("op", op)).Info("got role", slog.Any("role", ret))

	return ret.Role, nil
}
