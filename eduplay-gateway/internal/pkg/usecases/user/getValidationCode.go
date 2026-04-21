package user

import (
	"context"
	users "eduplay-gateway/internal/generated/clients/user"
	userModel "eduplay-gateway/internal/lib/models/user"
	errs "eduplay-gateway/internal/storage"
	"fmt"
	"log/slog"
)

func (a *UseCase) GetValidationCode(ctx context.Context, code userModel.Code) (bool, error) {
	const op = "Users.GetValidationCode"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to get validation code")

	message, err := a.userClient.GetVerificationCode(ctx, &users.Id{Id: code.Code})
	if err != nil {
		a.l.Error("failed to get validation code", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if message.Message == errs.ErrUserNotFound.Error() {
		return false, errs.ErrUserNotFound
	}
	if message.Message == errs.ErrCodeExpired.Error() {
		return false, errs.ErrCodeExpired
	}

	return true, nil
}
