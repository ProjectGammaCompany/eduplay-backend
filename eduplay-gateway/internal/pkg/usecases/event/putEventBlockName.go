package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutEventBlockName(ctx context.Context, req *eventModel.EventBlockName) (string, error) {
	const op = "event.UseCase.PutEventBlockName"

	s.log.With(slog.String("op", op)).Info("attempting to update event block name")

	var blockDto = &dto.Tag{
		Id:   req.BlockId,
		Name: req.Name,
	}

	ret, err := s.eventClient.PutEventBlockName(ctx, blockDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to update event block name", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("event block name updated", slog.Any("event", ret))

	return ret.Message, nil
}
