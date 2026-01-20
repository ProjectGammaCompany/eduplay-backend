package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error) {
	const op = "Events.UseCase.PostEventBlock"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("creating event block")

	id, err := a.storage.PostEventBlock(ctx, in)
	if err != nil {
		log.Error("failed to create event block", err.Error(), slog.String("event", in.EventId))
		return "", err
	}

	return id, nil
}
