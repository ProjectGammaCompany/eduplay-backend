package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutEvent(ctx context.Context, in *dto.PutEventIn) (*dto.GetGroupsOut, error) {
	const op = "Events.UseCase.PutEvent"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating collaborator list")

	err := a.storage.UpdateEventCollaborators(ctx, in.EventId, in.Collaborators)
	if err != nil {
		log.Error("failed to update collaborator list", err.Error(), slog.String("event", in.EventId))
		return nil, err
	}

	log.Info("updating groups")

	err = a.storage.UpdateEventGroups(ctx, in.EventId, in.Groups)
	if err != nil {
		log.Error("failed to update groups", err.Error(), slog.String("event", in.EventId))
		return nil, err
	}

	log.Info("updating event")

	_, err = a.storage.PutEvent(ctx, in)
	if err != nil {
		log.Error("failed to update event", err.Error(), slog.String("event", in.EventId))
		return nil, err
	}

	log.Info("event updated")

	groups, err := a.storage.GetGroups(ctx, in.EventId)
	if err != nil {
		log.Error("failed to get groups in event", err.Error(), slog.String("event", in.EventId))
		return nil, err
	}

	return groups, nil
}
