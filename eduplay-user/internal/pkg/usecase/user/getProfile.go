package user

import (
	"context"
	"log/slog"

	dto "eduplay-user/internal/generated"
)

func (a *UseCase) GetProfile(ctx context.Context, userId string) (*dto.Profile, error) {
	const op = "Users.GetProfile"

	log := a.log.With(
		slog.String("op", op),
	)

	profile, err := a.storage.GetProfile(ctx, userId)

	if err != nil {
		return nil, err
	}

	log.Info("got user profile successfully")

	return profile, nil
}
