package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetPublicEvent(ctx context.Context, in *dto.UserEventIds) (*dto.GetPublicEvent, error) {
	const op = "Events.UseCase.GetPublicEvent"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting public event")

	events, err := a.storage.GetPublicEvent(ctx, in)
	if err != nil {
		log.Error("failed to get public event", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	return events, nil
}
