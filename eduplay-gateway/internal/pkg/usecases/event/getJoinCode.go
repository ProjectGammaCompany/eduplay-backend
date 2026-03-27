package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetJoinCode(ctx context.Context, eventId string) (*eventModel.JoinCode, error) {
	const op = "event.UseCase.GetJoinCode"

	s.log.With(slog.String("op", op)).Info("attempting to get join code")

	ret, err := s.eventClient.GetJoinCode(ctx, &eventDto.Id{Id: eventId})

	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get join code", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("complaint posted", slog.Any("event", ret))

	return &eventModel.JoinCode{
		EventId:   ret.EventId,
		JoinCode:  ret.JoinCode,
		ExpiresAt: ret.ExpiresAt.AsTime().Format("02.01.2006 15:04:05.000"),
	}, nil
}
