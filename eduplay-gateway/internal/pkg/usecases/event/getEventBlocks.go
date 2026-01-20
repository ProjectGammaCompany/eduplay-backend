package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetEventBlocks(ctx context.Context, req *eventModel.Id) (*eventModel.GetEventBlocksOut, error) {
	const op = "event.UseCase.GetEventBlocks"

	s.log.With(slog.String("op", op)).Info("attempting to get event blocks")

	ret, err := s.eventClient.GetEventBlocks(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", err.Error())
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got event", slog.Any("event", ret))

	return eventModel.GetEventBlocksFromDto(ret), nil
}
