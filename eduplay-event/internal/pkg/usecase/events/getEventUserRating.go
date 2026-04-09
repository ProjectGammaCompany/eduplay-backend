package event

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	dto "eduplay-event/internal/generated"
	errs "eduplay-event/internal/storage"
)

func (a *UseCase) GetEventUserRating(ctx context.Context, in *dto.UserEventIds) (*dto.MessageOut, error) {
	const op = "Events.UseCase.GetEventUserRating"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting event user rating")

	rating, err := a.storage.GetEventUserRating(ctx, in.UserId, in.EventId)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &dto.MessageOut{Message: "-1"}, nil
		}
		log.Error("failed to get event user rating", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return &dto.MessageOut{Message: strconv.FormatInt(rating, 10)}, nil
}
