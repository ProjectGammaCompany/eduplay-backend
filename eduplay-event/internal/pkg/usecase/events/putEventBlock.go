package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error) {
	const op = "Events.UseCase.PutEventBlock"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating event block")

	id, err := a.storage.PutEventBlock(ctx, in)
	if err != nil {
		log.Error("failed to update event block", err.Error(), slog.String("event", in.EventId))
		return "", err
	}

	return id, nil
}
