package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostBlockCondition(ctx context.Context, in *dto.Condition) (*dto.PostConditionOut, error) {
	const op = "Events.UseCase.PostBlockCondition"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("creating condition")

	id, err := a.storage.PostBlockCondition(ctx, in)
	if err != nil {
		log.Error("failed to create condition", err.Error(), slog.String("block", in.PreviousBlockId))
		return nil, err
	}

	return id, nil
}
