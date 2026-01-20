package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetEvent(ctx context.Context, in *dto.Id) (*dto.PostEventIn, error) {
	const op = "Events.UseCase.GetEvent"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting event")

	event, err := a.storage.GetEvent(ctx, in.Id)
	if err != nil {
		log.Error("failed to get event", err.Error(), slog.String("event", in.Id))
		return nil, err
	}

	return event, nil
}
