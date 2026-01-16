package user

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/user"
	model "eduplay-gateway/internal/lib/models/user"
	"eduplay-gateway/internal/storage"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/metadata"
)

func (a *UseCase) Refresh(ctx context.Context, tokens *model.RefreshToken) (*model.RefreshToken, error) {
	const op = "Users.Refresh"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to refresh token")

	in := &dto.RefreshIn{
		RefreshToken: tokens.RefreshToken,
	}

	md := metadata.Pairs("Authorization", "Bearer "+tokens.AccessToken)

	newCtx := metadata.NewOutgoingContext(ctx, md)

	out, err := a.userClient.Refresh(newCtx, in)
	if err != nil {
		a.l.Error("failed to refresh token", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if out.Message != "" {
		log.Warn("refresh token refresh failed", slog.String("message", out.Message))
		if out.Message == "invalid refresh token" {
			return nil, storage.ErrInvalidRefreshToken
		}
		if out.Message == "refresh token expired" {
			return nil, storage.ErrRefreshTokenExpired
		}
		if out.Message == "refresh token not found in db" {
			return nil, storage.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	newTokens := &model.RefreshToken{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}

	return newTokens, nil
}
