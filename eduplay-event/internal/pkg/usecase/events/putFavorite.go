package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutFavorite(ctx context.Context, in *dto.PutFavoriteIn) (string, error) {
	const op = "Events.UseCase.PutFavorite"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("putting favorite")

	id, err := a.storage.PutFavorite(ctx, in)
	if err != nil {
		log.Error("failed to put favorite", err.Error(), slog.String("event", in.EventId))
		return "", err
	}

	return id, nil
}
