package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetBlockInfo(ctx context.Context, in *dto.Id) (*dto.PostEventBlockIn, error) {
	const op = "Events.UseCase.GetBlockInfo"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting block info")

	id, err := a.storage.GetBlockInfo(ctx, in.Id)
	if err != nil {
		log.Error("failed to get block info", err.Error(), slog.String("block", in.Id))
		return nil, err
	}

	return id, nil
}
