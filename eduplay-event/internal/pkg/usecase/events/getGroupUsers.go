package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetGroupUsers(ctx context.Context, in *dto.Id) (*dto.GetGroupUsersOut, error) {
	const op = "Events.UseCase.GetGroupUsers"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting group users")

	ret, err := a.storage.GetGroupUsers(ctx, in.Id)
	if err != nil {
		log.Error("failed to get group users", err.Error(), slog.String("group", in.Id))
		return nil, err
	}

	return ret, nil
}
