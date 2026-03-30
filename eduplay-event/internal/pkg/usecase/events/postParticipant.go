package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostParticipant(ctx context.Context, in *dto.PostParticipantIn) (*dto.MessageOut, error) {
	const op = "Events.UseCase.PostParticipant"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("posting participant info")

	message, err := a.storage.PostParticipant(ctx, in.UserId, in.EventId, in.GroupId)
	if err != nil {
		log.Error("failed to post participant", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return &dto.MessageOut{Message: message}, nil
}
