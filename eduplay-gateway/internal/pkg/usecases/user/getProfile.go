package user

import (
	"context"
	model "eduplay-gateway/internal/lib/models/user"
	"fmt"
	"log/slog"
)

func (a *UseCase) GetProfile(ctx context.Context, userId string) (*model.Profile, error) {
	const op = "Users.GetProfile"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to get user profile")

	profile, err := a.userClient.GetProfile(ctx, userId)
	if err != nil {
		a.l.Error("failed to get user profile", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Profile{Email: profile.Email, Avatar: profile.Avatar, UserName: profile.UserName}, nil
}
