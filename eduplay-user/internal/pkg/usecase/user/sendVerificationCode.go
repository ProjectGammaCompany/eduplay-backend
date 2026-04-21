package user

import (
	"context"
	dto "eduplay-user/internal/generated"
	errs "eduplay-user/internal/storage"
	"errors"
	"log/slog"
)

func (a *UseCase) SendVerificationCode(ctx context.Context, email string) (*dto.MessageOut, error) {
	const op = "Users.SendVerificationCode"

	log := a.log.With(
		slog.String("op", op),
	)

	code, err := GenerateJoinCode(verificationCodeSize)
	if err != nil {
		a.log.Error("failed to generate verification code", slog.String("error", err.Error()))
		return nil, err
	}

	err = a.storage.PutVerificationCode(ctx, email, code)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return &dto.MessageOut{Message: err.Error()}, nil
		}
		a.log.Error("failed to put verification code", slog.String("error", err.Error()))
		return nil, err
	}

	a.SendVerificationCodeEmail(email, code)

	log.Info("sent verification code successfully")

	return &dto.MessageOut{Message: "OK"}, nil
}
