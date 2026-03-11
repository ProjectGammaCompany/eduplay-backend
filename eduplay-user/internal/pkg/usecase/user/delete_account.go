package user

import (
	"context"
	"log/slog"
)

func (a *UseCase) DeleteUserAccount(ctx context.Context, userId string) error {
	const op = "Users.DeleteUserAccount"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to delete user account")

	ret, err := a.rabbitmq.SendDeleteAccountMessage(ctx, userId)
	if err != nil {
		log.Error("failed to send delete account message", err.Error(), slog.String("userId", userId))
		return err
	}

	log.Info(ret)
	// err := a.storage.DeleteAccount(ctx, userId)
	// if err != nil {
	// 	return err
	// }

	return nil
}
