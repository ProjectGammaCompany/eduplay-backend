package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostTask(ctx context.Context, in *dto.Task) (string, error) {
	const op = "Events.UseCase.PostTask"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("creating task")

	id, err := a.storage.PostTask(ctx, in)
	if err != nil {
		log.Error("failed to create task", err.Error(), slog.String("block", in.BlockId))
		return "", err
	}

	return id, nil
}
