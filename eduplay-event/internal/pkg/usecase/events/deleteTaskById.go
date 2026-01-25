package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) DeleteTaskById(ctx context.Context, in *dto.Id) (string, error) {
	const op = "Events.UseCase.DeleteTaskById"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("deleting task by id")

	message, err := a.storage.DeleteTaskById(ctx, in.Id)
	if err != nil {
		log.Error("failed to delete task by id", err.Error(), slog.String("task", in.Id))
		return "", err
	}

	return message, nil
}
