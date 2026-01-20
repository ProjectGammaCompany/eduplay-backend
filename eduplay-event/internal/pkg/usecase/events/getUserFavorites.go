package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "Events.UseCase.GetUserFavorites"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting user favorites")

	events, err := a.storage.GetUserFavorites(ctx, in)
	if err != nil {
		log.Error("failed to get user favorites", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	return events, nil
}
