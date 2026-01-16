package user

// import (
// 	"context"
// 	dto "eduplay-gateway/internal/generated/clients/user"
// 	model "eduplay-gateway/internal/lib/models/user"
// 	"fmt"
// 	"log/slog"
// )

// func (a *UseCase) GetUserAccess(ctx context.Context, accessToken string) (*model.UserAccess, error) {
// 	const op = "Users.GetUserAccess"

// 	log := a.l.With(
// 		slog.String("op", op),
// 	)

// 	log.Info("attempting to get user access")

// 	in := &dto.GetUserAccessIn{
// 		AccessToken: accessToken,
// 	}

// 	// md := metadata.Pairs("Authorization", "Bearer "+tokens.AccessToken)

// 	// newCtx := metadata.NewOutgoingContext(ctx, md)

// 	out, err := a.userClient.GetUserAccess(ctx, in)
// 	if err != nil {
// 		a.l.Error("failed to refresh token", slog.String("error", err.Error()))

// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	userAccess := &model.UserAccess{
// 		UserId:      out.UserId,
// 		Role:        out.Role.String(),
// 		AccessLevel: out.AccessLevel,
// 	}

// 	return userAccess, nil
// }
