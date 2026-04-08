package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutBlockCondition(ctx context.Context, req *eventModel.Condition) (int64, error) {
	const op = "event.UseCase.PutBlockCondition"

	s.log.With(slog.String("op", op)).Info("attempting to update block condition")

	taskDto := eventModel.ConditionToDto(req)

	ret, err := s.eventClient.PutBlockCondition(ctx, taskDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to update block condition", slog.String("error", err.Error()))
		return 0, err
	}

	blockInfo, err := s.eventClient.GetBlockInfo(ctx, &dto.Id{Id: ret.Message})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get block info", slog.String("error", err.Error()))
		return 0, err
	}

	s.log.With(slog.String("op", op)).Info("block condition updated", slog.Any("event", ret))

	return blockInfo.Order, nil
}
