package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetOwnedEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "Events.UseCase.GetOwnedEvents"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting user owned events")

	events, err := a.storage.GetOwnedEvents(ctx, in)
	if err != nil {
		log.Error("failed to get user owned events", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	return events, nil
}
