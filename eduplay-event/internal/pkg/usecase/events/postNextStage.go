package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutNextStage(ctx context.Context, in *dto.EventBlockTaskUserIds) (string, error) {
	const op = "Events.UseCase.PutNextStage"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("putting next stage")

	message, err := a.storage.PutNextStage(ctx, in)
	if err != nil {
		log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId))
		return "", err
	}

	return message, nil
}
