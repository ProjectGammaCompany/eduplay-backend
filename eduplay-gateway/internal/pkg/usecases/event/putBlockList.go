package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutBlockList(ctx context.Context, req *eventModel.PutBlockListIn) (string, error) {
	const op = "event.UseCase.PutBlockList"

	s.log.With(slog.String("op", op)).Info("attempting to update block list order")

	message, err := s.eventClient.PutBlockList(ctx, eventModel.PutBlockListInToDto(req))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to update block list order", slog.String("error", err.Error()))
		return "", err
	}

	return message.Message, nil
}
