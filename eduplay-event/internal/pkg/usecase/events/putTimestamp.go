package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutTimestamp(ctx context.Context, in *dto.PutTimestampIn) (string, error) {
	const op = "Events.UseCase.PutTimestamp"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("putting timestamp")

	message, err := a.storage.PutTimestamp(ctx, in.UserId, in.EventId, in.Timestamp)
	if err != nil {
		log.Error("failed to put timestamp", err.Error(), slog.String("event", in.EventId))
		return "", err
	}

	return message, nil
}
