package notification

import (
	"context"
	"log/slog"

	dto "eduplay-notification/internal/generated"
)

func (a *UseCase) GetNotifications(ctx context.Context, in *dto.Filters) (*dto.NotificationInfos, error) {
	const op = "Users.GetNotifications"

	log := a.log.With(
		slog.String("op", op),
	)

	message, err := a.storage.GetNotifications(ctx, in)
	if err != nil {
		return nil, err
	}

	log.Info("notifications got successfully")

	return message, nil
}
