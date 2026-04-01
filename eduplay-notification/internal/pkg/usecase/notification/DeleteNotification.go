package notification

import (
	"context"
	"log/slog"

	dto "eduplay-notification/internal/generated"
)

func (a *UseCase) DeleteNotification(ctx context.Context, in *dto.Ids) error {
	const op = "Users.DeleteNotification"

	log := a.log.With(
		slog.String("op", op),
	)

	err := a.storage.DeleteNotification(ctx, in)
	if err != nil {
		return err
	}

	log.Info("notifications deleted successfully")

	return nil
}
