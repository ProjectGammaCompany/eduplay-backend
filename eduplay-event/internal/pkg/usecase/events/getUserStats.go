package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetUserStats(ctx context.Context, in *dto.UserEventIds) (*dto.User, error) {
	const op = "Events.UseCase.GetUserStats"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting user stats")

	ret, err := a.storage.GetUserStats(ctx, in.UserId, in.EventId)
	if err != nil {
		log.Error("failed to get user stats in event", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return ret, nil
}
