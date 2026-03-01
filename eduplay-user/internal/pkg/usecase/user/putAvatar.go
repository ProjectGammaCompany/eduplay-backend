package user

import (
	"context"
	"log/slog"

	dto "eduplay-user/internal/generated"
)

func (a *UseCase) PutAvatar(ctx context.Context, in *dto.Profile) (string, error) {
	const op = "Users.PutAvatar"

	log := a.log.With(
		slog.String("op", op),
	)

	message, err := a.storage.PutAvatar(ctx, in)
	if err != nil {
		return "", err
	}

	log.Info("put avatar successfully")

	return message, nil
}
