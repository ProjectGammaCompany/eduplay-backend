package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutGroups(ctx context.Context, in *dto.PutListIn) (string, error) {
	const op = "Events.UseCase.PutGroup"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating event groups")

	id, err := a.storage.PutGroupsInCondition(ctx, in)
	if err != nil {
		log.Error("failed to update event groups", err.Error(), slog.String("conditionId", in.Id))
		return "", err
	}

	return id, nil
}
