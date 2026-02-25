package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetUserStatus(ctx context.Context, in *dto.UserEventIds) (*dto.MessageOut, error) {
	const op = "Events.UseCase.GetUserStatus"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting user status")

	status, err := a.storage.GetUserStatus(ctx, in.UserId, in.EventId)
	if err != nil {
		log.Error("failed to get user status", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	return status, nil
}
