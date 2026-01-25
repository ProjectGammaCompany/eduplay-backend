package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetBlockConditions(ctx context.Context, in *dto.Id) (*dto.BlockInfo, error) {
	const op = "Events.UseCase.GetBlockConditions"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting block conditions")

	id, err := a.storage.GetBlockConditionsFull(ctx, in.Id)
	if err != nil {
		log.Error("failed to get block conditions", err.Error(), slog.String("block", in.Id))
		return nil, err
	}

	return id, nil
}
