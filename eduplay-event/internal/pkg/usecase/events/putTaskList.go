package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutTaskList(ctx context.Context, in *dto.PutListIn) (string, error) {
	const op = "Events.UseCase.PutTaskList"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating task list order")

	id, err := a.storage.PutTaskList(ctx, in)
	if err != nil {
		log.Error("failed to update task list order", err.Error(), slog.String("blockId", in.Id))
		return "", err
	}

	return id, nil
}
