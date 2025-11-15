package users

import (
	"context"
	users "eduplay-gateway/internal/generated/clients/users"
	"log/slog"
)

type UsersClient interface {
	SignUp(ctx context.Context, in *users.SignUpIn) (*users.SignUpOut, error)
	SignIn(ctx context.Context, in *users.SignInIn) (*users.SignUpOut, error)
	Refresh(ctx context.Context, in *users.RefreshIn) (*users.RefreshOut, error)
	GetUserAccess(ctx context.Context, in *users.GetUserAccessIn) (*users.GetUserAccessOut, error)
	GetUserData(ctx context.Context, token string) (*users.GetUserInfoOut, error)
	ChangeUserData(ctx context.Context, in *users.ChangeUserInfoIn) (*users.GetUserInfoOut, error)
	DeleteAccount(ctx context.Context, token string) error
	ChangePassword(ctx context.Context, in *users.ChangePasswordIn) error
	SignOutUser(ctx context.Context, token string) error
}

type UseCase struct {
	l           *slog.Logger
	usersClient UsersClient
}

func New(
	l *slog.Logger,
	usersCl UsersClient,
) *UseCase {
	return &UseCase{
		l:           l,
		usersClient: usersCl,
	}
}
