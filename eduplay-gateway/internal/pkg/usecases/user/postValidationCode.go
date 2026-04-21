package user

import (
	"context"
	users "eduplay-gateway/internal/generated/clients/user"
	userModel "eduplay-gateway/internal/lib/models/user"
	errs "eduplay-gateway/internal/storage"
	"fmt"
	"log/slog"
)

func (a *UseCase) PostValidationCode(ctx context.Context, email *userModel.Email) error {
	const op = "Users.PostValidationCode"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to send validation code")

	message, err := a.userClient.SendVerificationCode(ctx, &users.Id{Id: email.Email})
	if err != nil {
		a.l.Error("failed to send validation code", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	if message.Message == errs.ErrUserNotFound.Error() {
		return errs.ErrUserNotFound
	}

	return nil
}
