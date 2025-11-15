package users

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/users"
	model "eduplay-gateway/internal/lib/models/user"
	"eduplay-gateway/internal/storage"
	"fmt"
	"log/slog"
)

func (a *UseCase) SignUp(ctx context.Context, pd *model.SignUpIn) (*model.Credentials, error) {
	const op = "Users.Sign_Up"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to get user pd")

	in := &dto.SignUpIn{
		Name:         pd.Name,
		Surname:      pd.Surname,
		Email:        pd.Email,
		Organization: pd.Organization,
		Password:     pd.Password,
		Phone:        pd.Phone,
	}
	tokens, err := a.usersClient.SignUp(ctx, in)
	if err != nil {
		a.l.Error("failed to sign up", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if tokens.ErrorMessage != "" {
		return nil, storage.ErrUserAlreadyExists
	}

	credentials := &model.Credentials{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Role:         model.Role(tokens.Role.String()),
		AccessLevel:  tokens.AccessLevel,
	}

	return credentials, nil
}
