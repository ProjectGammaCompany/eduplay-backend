package user

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/metadata"
)

func (a *UseCase) DeleteAccount(ctx context.Context, token string) error {
	const op = "Users.DeleteAccount"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to delete user")

	md := metadata.Pairs("Authorization", "Bearer "+token)

	newCtx := metadata.NewOutgoingContext(ctx, md)

	err := a.userClient.DeleteAccount(newCtx, token)
	if err != nil {
		a.l.Error("failed to delete user", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
