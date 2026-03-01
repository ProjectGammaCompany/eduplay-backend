package usersClient

import (
	"context"
	users "eduplay-gateway/internal/generated/clients/user"
	"fmt"
	"log/slog"
	"time"

	interlog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	interretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api users.UsersClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retries int) (*Client, error) {
	const op = "UsersClient.New"

	retriesOpts := []interretry.CallOption{
		interretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		interretry.WithMax(uint(retries)),
		interretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []interlog.Option{
		interlog.WithLogOnEvents(interlog.PayloadReceived, interlog.PayloadSent),
	}

	// cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithChainUnaryInterceptor(
	// 		interlog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
	// 		interretry.UnaryClientInterceptor(retriesOpts...),
	// 	))

	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interlog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			interretry.UnaryClientInterceptor(retriesOpts...),
		))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{api: users.NewUsersClient(cc), log: log}, nil
}

func (cl *Client) SignUp(ctx context.Context, in *users.SignUpIn) (*users.SignUpOut, error) {
	op := "SignUp.Client"
	out, err := cl.api.SignUp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) SignIn(ctx context.Context, in *users.SignInIn) (*users.SignUpOut, error) {
	op := "SignIn.Client"
	out, err := cl.api.SignIn(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) Refresh(ctx context.Context, in *users.RefreshIn) (*users.RefreshOut, error) {
	op := "Refresh.Client"
	out, err := cl.api.Refresh(ctx, in)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutAvatar(ctx context.Context, in *users.Profile) (*users.Empty, error) {
	op := "PutAvatar.Client"
	out, err := cl.api.PutAvatar(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

// func (cl *Client) GetUserAccess(ctx context.Context, in *users.GetUserAccessIn) (*users.GetUserAccessOut, error) {
// 	op := "GetUserAccess.Client"
// 	out, err := cl.api.GetUserAccess(ctx, in)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return out, nil
// }

// func (cl *Client) GetUserData(ctx context.Context, token string) (*users.GetUserInfoOut, error) {
// 	op := "GetUserData.Client"
// 	out, err := cl.api.GetUserInfo(ctx, &users.GetUserAccessIn{AccessToken: token})
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return out, nil
// }

// func (cl *Client) ChangeUserData(ctx context.Context, in *users.ChangeUserInfoIn) (*users.GetUserInfoOut, error) {
// 	op := "ChangeUserData.Client"
// 	out, err := cl.api.ChangeUserInfo(ctx, in)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return out, nil
// }

func (cl *Client) DeleteAccount(ctx context.Context, token string) error {
	op := "DeleteAccount.Client"
	_, err := cl.api.DeleteAccount(ctx, &users.DeleteAccountIn{AccessToken: token})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (cl *Client) ChangePassword(ctx context.Context, in *users.ChangePasswordIn) error {
	op := "ChangePassword.Client"
	_, err := cl.api.ChangePassword(ctx, in)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (cl *Client) SignOutUser(ctx context.Context, token string) error {
	op := "SignOutUser.Client"
	_, err := cl.api.SignOut(ctx, &users.DeleteAccountIn{AccessToken: token})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (cl *Client) GetProfile(ctx context.Context, userId string) (*users.Profile, error) {
	op := "GetProfile.Client"
	out, err := cl.api.GetProfile(ctx, &users.Id{Id: userId})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func InterceptorLogger(l *slog.Logger) interlog.Logger {
	return interlog.LoggerFunc(func(ctx context.Context, lvl interlog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
