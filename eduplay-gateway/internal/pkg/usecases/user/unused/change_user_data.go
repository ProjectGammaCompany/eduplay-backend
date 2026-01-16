package user

// import (
// 	"context"
// 	users "eduplay-gateway/internal/generated/clients/user"
// 	userModel "eduplay-gateway/internal/lib/models/user"
// 	"fmt"
// 	"log/slog"

// 	"google.golang.org/grpc/metadata"
// )

// func (a *UseCase) UpdateUserData(ctx context.Context, info userModel.ChangeUserData, accessToken string) (*userModel.UserInfo, error) {
// 	const op = "Users.UpdateUserData"

// 	log := a.l.With(
// 		slog.String("op", op),
// 	)

// 	log.Info("attempting to change user data")

// 	md := metadata.Pairs("Authorization", "Bearer "+accessToken)

// 	newCtx := metadata.NewOutgoingContext(ctx, md)

// 	req := &users.ChangeUserInfoIn{
// 		Name:                   info.Name,
// 		Surname:                info.Surname,
// 		JobTitle:               info.JobTitle,
// 		Organisation:           info.Organisation,
// 		Email:                  info.Email,
// 		Phone:                  info.Phone,
// 		City:                   info.City,
// 		ShortOrganisationTitle: info.ShortOrganisationTitle,
// 		INN:                    info.INN,
// 		OrganisationType:       info.OrganisationType,
// 		AccessToken:            accessToken,
// 	}

// 	out, err := a.userClient.ChangeUserData(newCtx, req)
// 	if err != nil {
// 		a.l.Error("failed to change user data", slog.String("error", err.Error()))

// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	output := &userModel.UserInfo{
// 		Name:                   out.Name,
// 		Surname:                out.Surname,
// 		JobTitle:               out.JobTitle,
// 		Organisation:           out.Organisation,
// 		Phone:                  out.Phone,
// 		Email:                  out.Email,
// 		City:                   out.City,
// 		ShortOrganisationTitle: out.ShortOrganisationTitle,
// 		INN:                    out.INN,
// 		OrganisationType:       out.OrganisationType,
// 		CurrentTarrif:          out.CurrentTarrif,
// 	}

// 	return output, nil
// }
