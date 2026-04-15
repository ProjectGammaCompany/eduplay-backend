package user

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/user"
	eventModel "eduplay-gateway/internal/lib/models/user"
	"log/slog"
)

func (s *UseCase) PutUsername(ctx context.Context, req *eventModel.Profile) (string, error) {
	const op = "event.UseCase.PutUsername"

	s.l.With(slog.String("op", op)).Info("attempting to update username")

	var usernameDto = &dto.Profile{
		Email:    req.Email,
		UserName: req.UserName,
	}

	message, err := s.userClient.PutUsername(ctx, usernameDto)
	if err != nil {
		s.l.With(slog.String("op", op)).Error("failed to update username", slog.String("error", err.Error()))
		return "", err
	}

	s.l.With(slog.String("op", op)).Info("username updated", slog.Any("service answer", message))

	return "username updated", nil
}
