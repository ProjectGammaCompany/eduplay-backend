package event

import (
	"context"
	"log/slog"

	// "eduplay-event/internal/model"
	dto "eduplay-event/internal/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type storage interface {
	SaveFile(ctx context.Context, fileName string, fileUUID string) (string, error)
	PostEvent(ctx context.Context, in *dto.PostEventIn) (string, error)
	PutEvent(ctx context.Context, in *dto.PutEventIn) (string, error)
	GetEvent(ctx context.Context, id string) (*dto.PostEventIn, error)
	GetRole(ctx context.Context, userId string, eventId string) (int64, error)
	GetGroups(ctx context.Context, eventId string) (*dto.GetGroupsOut, error)
	GetCollaborators(ctx context.Context, eventId string) (*dto.GetCollaboratorsOut, error)
	PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	PutEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	GetEventBlocks(ctx context.Context, eventId string) (*dto.GetEventBlocksOut, error)
	GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetOwnedEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetHistory(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	PutFavorite(ctx context.Context, in *dto.PutFavoriteIn) (string, error)
	GetAllTags(ctx context.Context) (*dto.Tags, error)
	PostTask(ctx context.Context, in *dto.Task) (string, error)
	PutTask(ctx context.Context, in *dto.Task) (*dto.PutTaskOut, error)
	PostBlockCondition(ctx context.Context, in *dto.Condition) (*dto.PostConditionOut, error)
	PutBlockCondition(ctx context.Context, in *dto.Condition) (string, error)
	DeleteBlockCondition(ctx context.Context, conditionId string) (string, error)
	GetBlockInfo(ctx context.Context, id string) (*dto.PostEventBlockIn, error)
	GetBlockConditionsFull(ctx context.Context, id string) (*dto.BlockInfo, error)
	GetBlockTasks(ctx context.Context, blockId string) (*dto.Tasks, error)
	GetTaskById(ctx context.Context, taskId string) (*dto.Task, error)
	DeleteTaskById(ctx context.Context, taskId string) (string, error)
	DeleteEventBlock(ctx context.Context, blockId string) (string, error)
	DeleteEvent(ctx context.Context, eventId string) (string, error)
	PostAnswer(ctx context.Context, answer *dto.Answer) (string, error)
	GetPublicEvent(ctx context.Context, ids *dto.UserEventIds) (*dto.GetPublicEvent, error)
	GetNextStage(ctx context.Context, stage *dto.UserEventIds) (linkId string, currTaskId string, currBlockId string, finished bool, startTime *timestamppb.Timestamp, err error)
	PutNextStage(ctx context.Context, stage *dto.EventBlockTaskUserIds) (string, error)
	PutTimestamp(ctx context.Context, userId string, eventId string, timestamp *timestamppb.Timestamp) (string, error)
	EndMe(ctx context.Context, userId string, eventId string) (string, error)
	GetUserBlockPointsSum(ctx context.Context, userId string, blockId string) (int64, error)
	GetUserBlockTasksShort(ctx context.Context, blockId string, userId string) ([]*dto.NextStageTaskShort, error)
	GetUserStatus(ctx context.Context, userId string, eventId string) (*dto.MessageOut, error)
	GetCollaboratorIds(ctx context.Context, emails []string) ([]string, error)
	UpdateEventCollaborators(ctx context.Context, eventId string, collaboratorIds []string) error
	UpdateEventGroups(ctx context.Context, eventId string, groups []*dto.Group) error
}

type UseCase struct {
	log     *slog.Logger
	storage storage
	secret  string
}

func New(
	log *slog.Logger,
	st storage,
	secret string,
) *UseCase {
	return &UseCase{
		log:     log,
		storage: st,
		secret:  secret,
	}
}
