package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetBlockTasks(ctx context.Context, in *dto.Id) (*dto.Tasks, error) {
	const op = "Events.UseCase.GetBlockTasks"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting block tasks")

	tasks, err := a.storage.GetBlockTasks(ctx, in.Id)
	if err != nil {
		log.Error("failed to get block tasks", err.Error(), slog.String("block", in.Id))
		return nil, err
	}

	return tasks, nil
}
