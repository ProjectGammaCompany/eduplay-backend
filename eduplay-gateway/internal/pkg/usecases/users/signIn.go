package users

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/users"
	model "eduplay-gateway/internal/lib/models/user"
	"eduplay-gateway/internal/storage"
	"fmt"
	"log/slog"
)

func (a *UseCase) SignIn(ctx context.Context, pd *model.UserPD) (*model.Credentials, error) {
	const op = "Users.Sign_In"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to get user pd")

	in := &dto.SignInIn{
		Email:    pd.Email,
		Password: pd.Password,
	}
	tokens, err := a.usersClient.SignIn(ctx, in)
	if err != nil {
		a.l.Error("failed to sign in", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if tokens.ErrorMessage == "user not found" {
		return nil, storage.ErrUserNotFound
	}

	if tokens.ErrorMessage == "incorrect password" {
		return nil, storage.ErrIncorrectPassword
	}

	if tokens.ErrorMessage == "user is already active" {
		return nil, storage.ErrIsActive
	}

	credentials := &model.Credentials{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Role:         model.Role(tokens.Role.String()),
		AccessLevel:  tokens.AccessLevel,
	}

	return credentials, nil
}
