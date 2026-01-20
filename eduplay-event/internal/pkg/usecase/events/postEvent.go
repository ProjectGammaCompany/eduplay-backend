package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostEvent(ctx context.Context, in *dto.PostEventIn) (string, error) {
	const op = "Events.UseCase.PostEvent"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("creating event")

	id, err := a.storage.PostEvent(ctx, in)
	if err != nil {
		log.Error("failed to create event", err.Error(), slog.String("event", in.Title))
		return "", err
	}

	return id, nil
}
