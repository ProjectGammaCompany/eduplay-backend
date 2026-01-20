package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetRole(ctx context.Context, in *dto.GetRoleIn) (*dto.GetRoleOut, error) {
	const op = "Events.UseCase.GetRole"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting role")

	num, err := a.storage.GetRole(ctx, in.UserId, in.EventId)
	if err != nil {
		log.Error("failed to get user role in event", err.Error(), slog.String("event", in.EventId))
		return nil, err
	}

	return &dto.GetRoleOut{Role: num}, nil
}
