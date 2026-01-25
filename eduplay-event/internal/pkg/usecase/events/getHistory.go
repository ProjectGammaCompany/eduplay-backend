package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetHistory(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "Events.UseCase.GetHistory"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting user history")

	events, err := a.storage.GetHistory(ctx, in)
	if err != nil {
		log.Error("failed to get user history", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	return events, nil
}
