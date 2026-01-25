package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) DeleteEventById(ctx context.Context, in *dto.Id) (string, error) {
	const op = "Events.UseCase.DeleteEventById"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("deleting event by id")

	message, err := a.storage.DeleteEvent(ctx, in.Id)
	if err != nil {
		log.Error("failed to delete event by id", err.Error(), slog.String("event", in.Id))
		return "", err
	}

	return message, nil
}
