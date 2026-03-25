package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetEventUsers(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error) {
	const op = "Events.UseCase.GetEventUsers"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting event users")

	ret, err := a.storage.GetEventUsers(ctx, in.Id)
	if err != nil {
		log.Error("failed to get event users", err.Error(), slog.String("event", in.Id))
		return nil, err
	}

	return ret, nil
}
