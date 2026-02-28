package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PutBlockCondition(ctx context.Context, in *dto.Condition) (*dto.MessageOut, error) {
	const op = "Events.UseCase.PutBlockCondition"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating condition")

	message, err := a.storage.PutBlockCondition(ctx, in)
	if err != nil {
		log.Error("failed to update condition", err.Error(), slog.String("block", in.PreviousBlockId))
		return nil, err
	}

	return &dto.MessageOut{Message: message}, nil
}
