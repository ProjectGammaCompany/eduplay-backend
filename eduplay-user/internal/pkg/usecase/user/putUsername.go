package user

import (
	"context"
	"log/slog"

	dto "eduplay-user/internal/generated"
)

func (a *UseCase) PutUsername(ctx context.Context, in *dto.Profile) (string, error) {
	const op = "Users.PutUsername"

	log := a.log.With(
		slog.String("op", op),
	)

	message, err := a.storage.PutUsername(ctx, in)
	if err != nil {
		return "", err
	}

	log.Info("put username successfully")

	return message, nil
}
