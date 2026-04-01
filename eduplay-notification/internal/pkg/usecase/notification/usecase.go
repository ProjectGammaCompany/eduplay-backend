package notification

import (
	"context"
	dto "eduplay-notification/internal/generated"
	"log/slog"
)

type storage interface {
	GetNotifications(ctx context.Context, in *dto.Filters) (*dto.NotificationInfos, error)
	DeleteNotification(ctx context.Context, in *dto.Ids) error
}

type UseCase struct {
	log     *slog.Logger
	storage storage
	secret  string
}

func New(
	log *slog.Logger,
	st storage,
	secret string,
) *UseCase {
	return &UseCase{
		log:     log,
		storage: st,
		secret:  secret,
	}
}
