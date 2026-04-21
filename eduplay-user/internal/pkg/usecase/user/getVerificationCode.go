package user

import (
	"context"
	dto "eduplay-user/internal/generated"
	errs "eduplay-user/internal/storage"
	"errors"
	"log/slog"
)

func (a *UseCase) GetVerificationCode(ctx context.Context, code string) (*dto.MessageOut, error) {
	const op = "Users.GetVerificationCode"

	log := a.log.With(
		slog.String("op", op),
	)

	_, _, err := a.storage.GetVerificationCode(ctx, code)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) || errors.Is(err, errs.ErrCodeExpired) {
			return &dto.MessageOut{Message: err.Error()}, nil
		}
		a.log.Error("failed to get verification code", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("got verification code successfully")

	return &dto.MessageOut{Message: "OK"}, nil
}
