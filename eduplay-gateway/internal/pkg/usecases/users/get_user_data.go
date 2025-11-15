package users

import (
	"context"
	model "eduplay-gateway/internal/lib/models/user"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/metadata"
)

func (a *UseCase) GetUserData(ctx context.Context, token string) (*model.UserInfo, error) {
	const op = "Users.GetUserData"

	log := a.l.With(
		slog.String("op", op),
	)

	log.Info("attempting to get user data")

	md := metadata.Pairs("Authorization", "Bearer "+token)

	newCtx := metadata.NewOutgoingContext(ctx, md)

	out, err := a.usersClient.GetUserData(newCtx, token)
	if err != nil {
		a.l.Error("failed to get user data", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	info := &model.UserInfo{
		Name:                   out.Name,
		Surname:                out.Surname,
		JobTitle:               out.JobTitle,
		Organisation:           out.Organisation,
		Phone:                  out.Phone,
		Email:                  out.Email,
		City:                   out.City,
		ShortOrganisationTitle: out.ShortOrganisationTitle,
		INN:                    out.INN,
		OrganisationType:       out.OrganisationType,
		CurrentTarrif:          out.CurrentTarrif,
	}

	return info, nil
}
