package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetBlocksForConditions(ctx context.Context, req *eventModel.Id) (*eventModel.GetBlocksForConditionsOut, error) {
	const op = "event.UseCase.GetBlocksForConditions"

	s.log.With(slog.String("op", op)).Info("attempting to get event blocks for conditions")

	ret, err := s.eventClient.GetEventBlocks(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event blocks for conditions", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got event", slog.Any("event", ret))

	return eventModel.GetBlocksForConditionsOutFromDto(ret), nil
}
