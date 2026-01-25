package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) DeleteBlockCondition(ctx context.Context, in *dto.Id) (string, error) {
	const op = "Events.UseCase.DeleteBlockCondition"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("deleting block condition by id")

	message, err := a.storage.DeleteBlockCondition(ctx, in.Id)
	if err != nil {
		log.Error("failed to delete block condition by id", err.Error(), slog.String("block condition", in.Id))
		return "", err
	}

	return message, nil
}
