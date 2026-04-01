package notification

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/notification"
	model "eduplay-gateway/internal/lib/models/notification"
	"fmt"
	"log/slog"
)

func (a *UseCase) GetNotifications(ctx context.Context, in *model.NotificationFilter) (*model.Notifications, error) {
	const op = "Notification.GetNotifications"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to get notifications")

	notifications, err := a.notifClient.GetNotifications(ctx, &dto.Filters{UserId: in.UserId, Page: in.Page, MaxOnPage: in.MaxOnPage})
	if err != nil {
		a.l.Error("failed to get notifications", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return model.NotificationsFromDto(notifications), nil
}
