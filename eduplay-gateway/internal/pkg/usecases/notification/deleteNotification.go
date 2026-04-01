package notification

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/notification"
	"fmt"
	"log/slog"
)

func (a *UseCase) DeleteNotification(ctx context.Context, userId string, notificationId string) error {
	const op = "Notification.DeleteNotification"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to delete notifications")

	err := a.notifClient.DeleteNotification(ctx, &dto.Ids{UserId: userId, NotificationId: notificationId})
	if err != nil {
		a.l.Error("failed to delete notifications", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
