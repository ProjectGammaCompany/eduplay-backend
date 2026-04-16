package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"eduplay-gateway/internal/storage"
	"log/slog"
)

func (s *UseCase) PostRate(ctx context.Context, req *eventModel.Rate) (string, error) {
	const op = "event.UseCase.PostRate"

	s.log.With(slog.String("op", op)).Info("attempting to post user rate")

	currStatus, err := s.eventClient.GetUserStatus(ctx, &dto.UserEventIds{UserId: req.UserId, EventId: req.EventId})

	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get user status", slog.String("error", err.Error()))
		return "", err
	}

	if currStatus.Status != "finished" {
		s.log.With(slog.String("op", op)).Error("user has not finished event", slog.String("event", req.EventId), slog.String("user", req.UserId))
		return "", storage.ErrInvalidOperation
	}

	ret, err := s.eventClient.PostRate(ctx, &dto.Rate{
		EventId: req.EventId,
		UserId:  req.UserId,
		Rate:    req.Rate,
	})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post user rate", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("user rate posted", slog.Any("event", req.EventId), slog.Any("user", req.UserId), slog.Any("rate", req.Rate))

	return ret.Message, nil
}
