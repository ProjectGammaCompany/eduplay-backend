package notificationClient

import (
	"context"
	notification "eduplay-gateway/internal/generated/clients/notification"
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
	api notification.NotificationsClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retries int) (*Client, error) {
	const op = "NotificationClient.New"

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

	return &Client{api: notification.NewNotificationsClient(cc), log: log}, nil
}

func (cl *Client) GetNotifications(ctx context.Context, in *notification.Filters) (*notification.NotificationInfos, error) {
	op := "GetNotifications.Client"

	out, err := cl.api.GetNotifications(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) DeleteNotification(ctx context.Context, in *notification.Ids) error {
	op := "DeleteNotification.Client"
	_, err := cl.api.DeleteNotification(ctx, in)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func InterceptorLogger(l *slog.Logger) interlog.Logger {
	return interlog.LoggerFunc(func(ctx context.Context, lvl interlog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
