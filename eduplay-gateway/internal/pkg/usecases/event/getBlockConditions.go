package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetBlockConditions(ctx context.Context, req *eventModel.Id) (*eventModel.Conditions, error) {
	const op = "event.UseCase.GetBlockInfo"

	s.log.With(slog.String("op", op)).Info("attempting to ge block Conditions")

	ret, err := s.eventClient.GetBlockConditions(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get block conditions", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got block conditions", slog.Any("event", ret))

	return eventModel.ConditionsFromDto(ret.Conditions), nil
}
