package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetCollaborators(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error) {
	const op = "Events.UseCase.GetCollaborators"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting collaborators")

	groups, err := a.storage.GetCollaborators(ctx, in.Id)
	if err != nil {
		log.Error("failed to get collaborators in event", err.Error(), slog.String("event", in.Id))
		return nil, err
	}

	return groups, nil
}
