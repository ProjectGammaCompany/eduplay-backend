package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetGroups(ctx context.Context, in *dto.Id) (*dto.GetGroupsOut, error) {
	const op = "Events.UseCase.GetGroups"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting groups")

	groups, err := a.storage.GetGroups(ctx, in.Id)
	if err != nil {
		log.Error("failed to get groups in event", err.Error(), slog.String("event", in.Id))
		return nil, err
	}

	return groups, nil
}
