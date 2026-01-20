package eventsClient

import (
	"context"
	events "eduplay-gateway/internal/generated/clients/event"
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
	api events.EventsClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retries int) (*Client, error) {
	const op = "EventsClient.New"

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

	return &Client{api: events.NewEventsClient(cc), log: log}, nil
}

func (cl *Client) SaveFile(ctx context.Context, in *events.SaveFileIn) (*events.MessageOut, error) {
	op := "SaveFile.Client"
	out, err := cl.api.SaveFile(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PostEvent(ctx context.Context, in *events.PostEventIn) (*events.MessageOut, error) {
	op := "PostEvent.Client"
	out, err := cl.api.PostEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetEvent(ctx context.Context, in *events.Id) (*events.PostEventIn, error) {
	op := "GetEvent.Client"
	out, err := cl.api.GetEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetRole(ctx context.Context, in *events.GetRoleIn) (*events.GetRoleOut, error) {
	op := "GetRole.Client"
	out, err := cl.api.GetRole(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetGroups(ctx context.Context, in *events.Id) (*events.GetGroupsOut, error) {
	op := "GetGroups.Client"
	out, err := cl.api.GetGroups(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetCollaborators(ctx context.Context, in *events.Id) (*events.GetCollaboratorsOut, error) {
	op := "GetCollaborators.Client"
	out, err := cl.api.GetCollaborators(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PostEventBlock(ctx context.Context, in *events.PostEventBlockIn) (*events.MessageOut, error) {
	op := "PostEventBlock.Client"
	out, err := cl.api.PostEventBlock(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetEventBlocks(ctx context.Context, in *events.Id) (*events.GetEventBlocksOut, error) {
	op := "GetEventBlocks.Client"
	out, err := cl.api.GetEventBlocks(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetPublicEvents(ctx context.Context, in *events.EventBaseFilters) (*events.GetPublicEventsOut, error) {
	op := "GetPublicEvents.Client"
	out, err := cl.api.GetPublicEvents(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetUserFavorites(ctx context.Context, in *events.EventBaseFilters) (*events.GetPublicEventsOut, error) {
	op := "GetUserFavorites.Client"
	out, err := cl.api.GetUserFavorites(ctx, in)
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
