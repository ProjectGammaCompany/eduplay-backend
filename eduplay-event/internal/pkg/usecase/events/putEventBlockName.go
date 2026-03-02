package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutEventBlockName(ctx context.Context, in *dto.Tag) (string, error) {
	const op = "Events.UseCase.PutEventBlockName"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating event block name")

	id, err := a.storage.PutEventBlockName(ctx, in)
	if err != nil {
		log.Error("failed to update event block name", err.Error(), slog.String("blockId", in.Id))
		return "", err
	}

	return id, nil
}
