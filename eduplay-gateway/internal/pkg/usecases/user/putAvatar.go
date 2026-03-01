package user

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/user"
	eventModel "eduplay-gateway/internal/lib/models/user"
	"log/slog"
)

func (s *UseCase) PutAvatar(ctx context.Context, req *eventModel.Profile) (string, error) {
	const op = "event.UseCase.PutAvatar"

	s.l.With(slog.String("op", op)).Info("attempting to update avatar")

	var avatarDto = &dto.Profile{
		Email:  req.Email,
		Avatar: req.Avatar,
	}

	message, err := s.userClient.PutAvatar(ctx, avatarDto)
	if err != nil {
		s.l.With(slog.String("op", op)).Error("failed to update avatar", slog.String("error", err.Error()))
		return "", err
	}

	s.l.With(slog.String("op", op)).Info("avatar updated", slog.Any("service answer", message))

	return "avatar updated", nil
}
