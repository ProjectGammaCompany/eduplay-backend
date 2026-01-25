package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) DeleteBlockById(ctx context.Context, in *dto.Id) (string, error) {
	const op = "Events.UseCase.DeleteBlockById"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("deleting block by id")

	message, err := a.storage.DeleteEventBlock(ctx, in.Id)
	if err != nil {
		log.Error("failed to delete block by id", err.Error(), slog.String("block", in.Id))
		return "", err
	}

	return message, nil
}
