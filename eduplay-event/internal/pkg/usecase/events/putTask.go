package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutTask(ctx context.Context, in *dto.Task) (*dto.PutTaskOut, error) {
	const op = "Events.UseCase.PutTask"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating task")

	ret, err := a.storage.PutTask(ctx, in)
	if err != nil {
		log.Error("failed to update task", err.Error(), slog.String("block", in.BlockId))
		return nil, err
	}

	return ret, nil
}
