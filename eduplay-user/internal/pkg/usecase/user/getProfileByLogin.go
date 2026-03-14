package user

import (
	"context"
	"log/slog"

	dto "eduplay-user/internal/generated"
)

func (a *UseCase) GetProfileByLogin(ctx context.Context, login string) (*dto.Profile, error) {
	const op = "Users.GetProfileByLogin"

	log := a.log.With(
		slog.String("op", op),
	)

	profile, err := a.storage.GetProfileByLogin(ctx, login)

	if err != nil {
		return nil, err
	}

	log.Info("got user profile by login successfully")

	return profile, nil
}
