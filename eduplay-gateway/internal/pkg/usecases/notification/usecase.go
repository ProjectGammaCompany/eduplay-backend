package notification

import (
	"context"
	notification "eduplay-gateway/internal/generated/clients/notification"
	"log/slog"
)

type NotificationClient interface {
	GetNotifications(ctx context.Context, in *notification.Filters) (*notification.NotificationInfos, error)
	DeleteNotification(ctx context.Context, in *notification.Ids) error
}

type UseCase struct {
	l           *slog.Logger
	notifClient NotificationClient
}

func New(
	l *slog.Logger,
	notifCl NotificationClient,
) *UseCase {
	return &UseCase{
		l:           l,
		notifClient: notifCl,
	}
}
