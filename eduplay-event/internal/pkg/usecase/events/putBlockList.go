package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutBlockList(ctx context.Context, in *dto.PutListIn) (string, error) {
	const op = "Events.UseCase.PutBlockList"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating block list order")

	id, err := a.storage.PutBlockList(ctx, in)
	if err != nil {
		log.Error("failed to update block list order", err.Error(), slog.String("blockId", in.Id))
		return "", err
	}

	return id, nil
}
