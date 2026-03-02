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

func (cl *Client) PutEvent(ctx context.Context, in *events.PutEventIn) (*events.GetGroupsOut, error) {
	op := "PutEvent.Client"
	out, err := cl.api.PutEvent(ctx, in)
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

func (cl *Client) PutGroups(ctx context.Context, in *events.PutListIn) (*events.MessageOut, error) {
	op := "PutGroups.Client"
	out, err := cl.api.PutGroups(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutTaskList(ctx context.Context, in *events.PutListIn) (*events.MessageOut, error) {
	op := "PutTaskList.Client"
	out, err := cl.api.PutTaskList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutBlockList(ctx context.Context, in *events.PutListIn) (*events.MessageOut, error) {
	op := "PutBlockList.Client"
	out, err := cl.api.PutBlockList(ctx, in)
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

func (cl *Client) PutEventBlock(ctx context.Context, in *events.PostEventBlockIn) (*events.MessageOut, error) {
	op := "PutEventBlock.Client"
	out, err := cl.api.PutEventBlock(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutEventBlockName(ctx context.Context, in *events.Tag) (*events.MessageOut, error) {
	op := "PutEventBlockName.Client"
	out, err := cl.api.PutEventBlockName(ctx, in)
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

func (cl *Client) GetOwnedEvents(ctx context.Context, in *events.EventBaseFilters) (*events.GetPublicEventsOut, error) {
	op := "GetOwnedEvents.Client"
	out, err := cl.api.GetOwnedEvents(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetHistory(ctx context.Context, in *events.EventBaseFilters) (*events.GetPublicEventsOut, error) {
	op := "GetHistoryEvents.Client"
	out, err := cl.api.GetHistory(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutFavorite(ctx context.Context, in *events.PutFavoriteIn) (*events.MessageOut, error) {
	op := "PutFavorite.Client"
	out, err := cl.api.PutFavorite(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetAllTags(ctx context.Context) (*events.Tags, error) {
	op := "GetAllTags.Client"
	out, err := cl.api.GetAllTags(ctx, &events.Empty{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PostTask(ctx context.Context, in *events.Task) (*events.MessageOut, error) {
	op := "PostTask.Client"
	out, err := cl.api.PostTask(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutTask(ctx context.Context, in *events.Task) (*events.PutTaskOut, error) {
	op := "PutTask.Client"
	out, err := cl.api.PutTask(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PostBlockCondition(ctx context.Context, in *events.Condition) (*events.PostConditionOut, error) {
	op := "PostBlockCondition.Client"
	out, err := cl.api.PostBlockCondition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutBlockCondition(ctx context.Context, in *events.Condition) (*events.MessageOut, error) {
	op := "PutBlockCondition.Client"
	out, err := cl.api.PutBlockCondition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetBlockInfo(ctx context.Context, in *events.Id) (*events.PostEventBlockIn, error) {
	op := "GetBlockInfo.Client"
	out, err := cl.api.GetBlockInfo(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetBlockConditions(ctx context.Context, in *events.Id) (*events.BlockInfo, error) {
	op := "GetBlockConditions.Client"
	out, err := cl.api.GetBlockConditions(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetBlockTasks(ctx context.Context, in *events.Id) (*events.Tasks, error) {
	op := "GetBlockTasks.Client"
	out, err := cl.api.GetBlockTasks(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetTaskById(ctx context.Context, in *events.Id) (*events.Task, error) {
	op := "GetTaskById.Client"
	out, err := cl.api.GetTaskById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) DeleteTaskById(ctx context.Context, in *events.Id) (*events.MessageOut, error) {
	op := "DeleteTaskById.Client"
	out, err := cl.api.DeleteTask(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PostAnswer(ctx context.Context, in *events.Answer) (*events.Answer, error) {
	op := "PostAnswer.Client"
	out, err := cl.api.PostAnswer(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) DeleteBlockById(ctx context.Context, in *events.Id) (*events.MessageOut, error) {
	op := "DeleteBlockById.Client"
	out, err := cl.api.DeleteBlockById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) DeleteEventById(ctx context.Context, in *events.Id) (*events.MessageOut, error) {
	op := "DeleteEventById.Storage"
	out, err := cl.api.DeleteEventById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) DeleteBlockCondition(ctx context.Context, in *events.Id) (*events.MessageOut, error) {
	op := "DeleteBlockCondition.Client"
	out, err := cl.api.DeleteBlockCondition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetEventForUser(ctx context.Context, in *events.UserEventIds) (*events.GetPublicEvent, error) {
	op := "GetEventForUser.Client"
	out, err := cl.api.GetEventForUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

// rpc PutNextStage(EventBlockTaskUserIds) returns (MessageOut);
// rpc GetNextStage(UserEventIds) returns (NextStageInfo);
// rpc PutTimestamp(PutTimestampIn) returns (MessageOut);

func (cl *Client) PutNextStage(ctx context.Context, in *events.EventBlockTaskUserIds) (*events.MessageOut, error) {
	op := "PutNextStage.Client"
	out, err := cl.api.PutNextStage(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetNextStage(ctx context.Context, in *events.UserEventIds) (*events.NextStageInfo, error) {
	op := "GetNextStage.Client"
	out, err := cl.api.GetNextStage(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) PutTimestamp(ctx context.Context, in *events.PutTimestampIn) (*events.MessageOut, error) {
	op := "PutTimestamp.Client"
	out, err := cl.api.PutTimestamp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (cl *Client) GetUserStatus(ctx context.Context, in *events.UserEventIds) (*events.MessageOut, error) {
	op := "GetUserStatus.Client"
	out, err := cl.api.GetUserStatus(ctx, in)
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
