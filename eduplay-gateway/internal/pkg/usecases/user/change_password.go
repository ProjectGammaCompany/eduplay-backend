package user

import (
	"context"
	users "eduplay-gateway/internal/generated/clients/user"
	userModel "eduplay-gateway/internal/lib/models/user"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/metadata"
)

func (a *UseCase) ChangePassword(ctx context.Context, info userModel.ChangePasswordRequest, accessToken string) error {
	const op = "Users.ChangePassword"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to change user password")

	md := metadata.Pairs("Authorization", "Bearer "+accessToken)

	newCtx := metadata.NewOutgoingContext(ctx, md)

	req := &users.ChangePasswordIn{
		Password:    info.Password,
		NewPassword: info.NewPassword,
		AccessToken: accessToken,
	}

	err := a.userClient.ChangePassword(newCtx, req)
	if err != nil {
		a.l.Error("failed to change user password", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
