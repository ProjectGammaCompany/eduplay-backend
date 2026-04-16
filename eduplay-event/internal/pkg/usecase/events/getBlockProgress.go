package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetBlockProgress(ctx context.Context, in *dto.UserEventIds) (*dto.BlockProgress, error) {
	const op = "Events.UseCase.GetBlockProgress"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting block progress")

	status, err := a.storage.GetBlockProgress(ctx, in)
	if err != nil {
		log.Error("failed to get block progress", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	return status, nil
}
