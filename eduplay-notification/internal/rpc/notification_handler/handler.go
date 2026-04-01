package notification_handler

import (
	"context"
	"fmt"

	dto "eduplay-notification/internal/generated"

	"google.golang.org/grpc"
)

type UseCase interface {
	GetNotifications(ctx context.Context, in *dto.Filters) (*dto.NotificationInfos, error)
	DeleteNotification(ctx context.Context, in *dto.Ids) error
}

type Handler struct {
	dto.UnimplementedNotificationsServer
	uc UseCase
}

func Register(gRPCServer *grpc.Server, uc UseCase) {
	dto.RegisterNotificationsServer(gRPCServer, &Handler{uc: uc})
}

func (h *Handler) GetNotifications(ctx context.Context, in *dto.Filters) (*dto.NotificationInfos, error) {
	op := "Notification.Handler.GetNotifications"

	notif, err := h.uc.GetNotifications(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return notif, nil
}

func (h *Handler) DeleteNotification(ctx context.Context, in *dto.Ids) (*dto.Empty, error) {
	op := "Notification.Handler.DeleteNotification"

	err := h.uc.DeleteNotification(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &dto.Empty{}, nil
}
