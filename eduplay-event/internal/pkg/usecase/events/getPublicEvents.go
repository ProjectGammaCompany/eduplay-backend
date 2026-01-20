package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "Events.UseCase.GetPublicEvents"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting public events")

	events, err := a.storage.GetPublicEvents(ctx, in)
	if err != nil {
		log.Error("failed to get public events", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	return events, nil
}
