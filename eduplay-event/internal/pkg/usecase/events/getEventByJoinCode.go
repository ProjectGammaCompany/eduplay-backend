package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetEventByJoinCode(ctx context.Context, in *dto.Id) (*dto.Id, error) {
	const op = "Events.UseCase.GetEventByJoinCode"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting event by join code")

	id, err := a.storage.GetEventByJoinCode(ctx, in.Id)
	if err != nil {
		log.Error("failed to get event by join code", err.Error(), slog.String("event", in.Id))
		return nil, err
	}

	return &dto.Id{Id: id}, nil
}
