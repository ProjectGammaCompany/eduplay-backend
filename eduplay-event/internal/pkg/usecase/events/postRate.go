package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostRate(ctx context.Context, in *dto.Rate) (*dto.MessageOut, error) {
	const op = "Events.UseCase.PostRate"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("posting user rate")

	message, err := a.storage.PostRate(ctx, in)
	if err != nil {
		log.Error("failed to post user rate", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId), slog.Any("rate", in.Rate))
		return nil, err
	}

	return message, nil
}
