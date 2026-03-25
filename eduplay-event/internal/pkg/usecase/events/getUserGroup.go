package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetUserGroup(ctx context.Context, in *dto.UserEventIds) (*dto.GetUserGroupOut, error) {
	const op = "Events.UseCase.GetUserGroup"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting user group")

	ret, err := a.storage.GetUserGroup(ctx, in.UserId, in.EventId)
	if err != nil {
		log.Error("failed to get user group in event", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return ret, nil
}
