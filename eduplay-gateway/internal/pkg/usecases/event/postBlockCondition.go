package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PostBlockCondition(ctx context.Context, req *eventModel.Condition) (*eventModel.PostConditionOut, error) {
	const op = "event.UseCase.PostBlockCondition"

	s.log.With(slog.String("op", op)).Info("attempting to post block condition")

	conditionDto := eventModel.ConditionToDto(req)

	ret, err := s.eventClient.PostBlockCondition(ctx, conditionDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post block condition", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("block condition posted", slog.Any("event", ret))

	return &eventModel.PostConditionOut{
		BlockOrder:  ret.BlockOrder,
		ConditionId: ret.ConditionId,
	}, nil
}
