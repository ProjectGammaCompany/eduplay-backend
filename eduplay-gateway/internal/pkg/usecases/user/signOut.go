package user

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/metadata"
)

func (a *UseCase) SignOutUser(ctx context.Context, token string) error {
	const op = "Users.SignOutUser"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to sign out user")

	md := metadata.Pairs("Authorization", "Bearer "+token)

	newCtx := metadata.NewOutgoingContext(ctx, md)

	err := a.userClient.SignOutUser(newCtx, token)
	if err != nil {
		a.l.Error("failed to sign out user", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
