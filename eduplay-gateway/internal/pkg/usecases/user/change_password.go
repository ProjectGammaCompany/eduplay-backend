package user

import (
	"context"
	users "eduplay-gateway/internal/generated/clients/user"
	model "eduplay-gateway/internal/lib/models/user"
	"fmt"
	"log/slog"
)

func (a *UseCase) ChangePassword(ctx context.Context, info model.ChangePasswordRequest) (*model.Credentials, error) {
	const op = "Users.ChangePassword"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to change user password")

	req := &users.ChangePasswordIn{
		Code:     info.Code,
		Password: info.Password,
	}

	session, err := a.userClient.ChangePassword(ctx, req)
	if err != nil {
		a.l.Error("failed to change user password", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Credentials{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}
