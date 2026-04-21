package user

import (
	"context"
	users "eduplay-gateway/internal/generated/clients/user"
	"log/slog"
)

type UserClient interface {
	SignUp(ctx context.Context, in *users.SignUpIn) (*users.SignUpOut, error)
	SignIn(ctx context.Context, in *users.SignInIn) (*users.SignUpOut, error)
	Refresh(ctx context.Context, in *users.RefreshIn) (*users.RefreshOut, error)
	PutAvatar(ctx context.Context, in *users.Profile) (*users.Empty, error)
	PutUsername(ctx context.Context, in *users.Profile) (*users.Empty, error)
	SendVerificationCode(ctx context.Context, in *users.Id) (*users.MessageOut, error)
	GetVerificationCode(ctx context.Context, in *users.Id) (*users.MessageOut, error)
	ChangePassword(ctx context.Context, in *users.ChangePasswordIn) (*users.SignUpOut, error)
	// GetUserAccess(ctx context.Context, in *users.GetUserAccessIn) (*users.GetUserAccessOut, error)
	// GetUserData(ctx context.Context, token string) (*users.GetUserInfoOut, error)
	// ChangeUserData(ctx context.Context, in *users.ChangeUserInfoIn) (*users.GetUserInfoOut, error)
	DeleteAccount(ctx context.Context, token string) error
	SignOutUser(ctx context.Context, token string) error
	GetProfile(ctx context.Context, userId string) (*users.Profile, error)
}

type UseCase struct {
	l          *slog.Logger
	userClient UserClient
}

func New(
	l *slog.Logger,
	userCl UserClient,
) *UseCase {
	return &UseCase{
		l:          l,
		userClient: userCl,
	}
}
