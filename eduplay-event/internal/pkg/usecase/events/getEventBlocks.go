package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetEventBlocks(ctx context.Context, in *dto.Id) (*dto.GetEventBlocksOut, error) {
	const op = "Events.UseCase.GetEventBlocks"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting event blocks")

	blocks, err := a.storage.GetEventBlocks(ctx, in.Id)
	if err != nil {
		log.Error("failed to get event blocks", err.Error(), slog.String("event", in.Id))
		return nil, err
	}

	return blocks, nil
}
