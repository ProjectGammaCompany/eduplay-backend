package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	userDto "eduplay-gateway/internal/generated/clients/user"
	"log/slog"
)

type EventClient interface {
	SaveFile(ctx context.Context, in *dto.SaveFileIn) (*dto.MessageOut, error)
	PostEvent(ctx context.Context, in *dto.PostEventIn) (*dto.MessageOut, error)
	PutEvent(ctx context.Context, in *dto.PutEventIn) (*dto.GetGroupsOut, error)
	GetEvent(ctx context.Context, in *dto.Id) (*dto.PostEventIn, error)
	GetRole(ctx context.Context, in *dto.GetRoleIn) (*dto.GetRoleOut, error)
	GetGroups(ctx context.Context, in *dto.Id) (*dto.GetGroupsOut, error)
	PutGroups(ctx context.Context, in *dto.PutListIn) (*dto.MessageOut, error)
	PutTaskList(ctx context.Context, in *dto.PutListIn) (*dto.MessageOut, error)
	PutBlockList(ctx context.Context, in *dto.PutListIn) (*dto.MessageOut, error)
	GetCollaborators(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error)
	PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (*dto.MessageOut, error)
	PutEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (*dto.MessageOut, error)
	PutEventBlockName(ctx context.Context, in *dto.Tag) (*dto.MessageOut, error)
	GetEventBlocks(ctx context.Context, in *dto.Id) (*dto.GetEventBlocksOut, error)
	GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetOwnedEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetHistory(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	PutFavorite(ctx context.Context, in *dto.PutFavoriteIn) (*dto.MessageOut, error)
	GetAllTags(ctx context.Context) (*dto.Tags, error)
	PostTask(ctx context.Context, in *dto.Task) (*dto.MessageOut, error)
	PutTask(ctx context.Context, in *dto.Task) (*dto.PutTaskOut, error)
	PostBlockCondition(ctx context.Context, in *dto.Condition) (*dto.PostConditionOut, error)
	PutBlockCondition(ctx context.Context, in *dto.Condition) (*dto.MessageOut, error)
	DeleteBlockCondition(ctx context.Context, in *dto.Id) (*dto.MessageOut, error)
	GetBlockInfo(ctx context.Context, in *dto.Id) (*dto.PostEventBlockIn, error)
	GetBlockConditions(ctx context.Context, in *dto.Id) (*dto.BlockInfo, error)
	GetBlockTasks(ctx context.Context, in *dto.Id) (*dto.Tasks, error)
	GetTaskById(ctx context.Context, in *dto.Id) (*dto.Task, error)
	DeleteTaskById(ctx context.Context, in *dto.Id) (*dto.MessageOut, error)
	DeleteBlockById(ctx context.Context, in *dto.Id) (*dto.MessageOut, error)
	DeleteEventById(ctx context.Context, in *dto.Id) (*dto.MessageOut, error)
	PostAnswer(ctx context.Context, in *dto.Answer) (*dto.Answer, error)
	GetEventForUser(ctx context.Context, in *dto.UserEventIds) (*dto.GetPublicEvent, error)
	PutNextStage(ctx context.Context, in *dto.EventBlockTaskUserIds) (*dto.MessageOut, error)
	GetNextStage(ctx context.Context, in *dto.UserEventIds) (*dto.NextStageInfo, error)
	PutTimestamp(ctx context.Context, in *dto.PutTimestampIn) (*dto.MessageOut, error)
	GetUserStatus(ctx context.Context, in *dto.UserEventIds) (*dto.MessageOut, error)
}

type UserClient interface {
	GetProfile(ctx context.Context, userId string) (*userDto.Profile, error)
}

type UseCase struct {
	eventClient EventClient
	userClient  UserClient
	log         *slog.Logger
}

func New(log *slog.Logger, eventClient EventClient, userClient UserClient) *UseCase {
	return &UseCase{log: log, eventClient: eventClient, userClient: userClient}
}
