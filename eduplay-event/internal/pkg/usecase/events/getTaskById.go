package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetTaskById(ctx context.Context, in *dto.Id) (*dto.Task, error) {
	const op = "Events.UseCase.GetTaskById"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting task by id")

	tasks, err := a.storage.GetTaskById(ctx, in.Id)
	if err != nil {
		log.Error("failed to get task by id", err.Error(), slog.String("block", in.Id))
		return nil, err
	}

	return tasks, nil
}
