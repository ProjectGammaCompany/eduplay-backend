package event

import (
	"context"
	"log/slog"
	"time"

	// "eduplay-event/internal/model"
	dto "eduplay-event/internal/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type storage interface {
	SaveFile(ctx context.Context, fileName string, fileKey string, fileUUID string) (string, error)
	PostEvent(ctx context.Context, in *dto.PostEventIn) (string, error)
	PutEvent(ctx context.Context, in *dto.PutEventIn) (string, error)
	GetEvent(ctx context.Context, id string) (*dto.PostEventIn, error)
	GetRole(ctx context.Context, userId string, eventId string) (int64, error)
	GetGroups(ctx context.Context, eventId string) (*dto.GetGroupsOut, error)
	PutGroupsInCondition(ctx context.Context, in *dto.PutListIn) (string, error)
	PutTaskList(ctx context.Context, in *dto.PutListIn) (string, error)
	PutBlockList(ctx context.Context, in *dto.PutListIn) (string, error)
	GetCollaborators(ctx context.Context, eventId string) (*dto.GetCollaboratorsOut, error)
	PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	PutEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	PutEventBlockName(ctx context.Context, in *dto.Tag) (string, error)
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
	GetUserStatus(ctx context.Context, userId string, eventId string) (*dto.UserStatus, error)
	UpdateEventCollaborators(ctx context.Context, eventId string, collaboratorIds []string) error
	UpdateEventGroups(ctx context.Context, eventId string, groups []*dto.Group) error
	GetTaskAnswer(ctx context.Context, taskId string, userId string) (*dto.Answer, error)
	InsertJoinCode(ctx context.Context, eventId string, joinCode string) (*time.Time, error)
	GetJoinCode(ctx context.Context, eventId string) (*dto.JoinCode, error)
	GetUserStats(ctx context.Context, userId string, eventId string) (*dto.User, error)
	GetUserGroup(ctx context.Context, userId string, eventId string) (*dto.GetUserGroupOut, error)
	GetGroupUsers(ctx context.Context, groupId string) (*dto.GetGroupUsersOut, error)
	GetEventUsers(ctx context.Context, eventId string) (*dto.GetCollaboratorsOut, error)
	PostComplaint(ctx context.Context, in *dto.PostComplaintIn) (string, error)
	GetEventByJoinCode(ctx context.Context, joinCode string) (string, error)
	GetEventUserRating(ctx context.Context, userId string, eventId string) (int64, error)
	PostParticipant(ctx context.Context, userId string, eventId string, groupId string) (string, error)
	ClearBlockAnswers(ctx context.Context, userId string, blockId string) error
	PostRate(ctx context.Context, in *dto.Rate) (*dto.MessageOut, error)
	GetBlockProgress(ctx context.Context, in *dto.UserEventIds) (*dto.BlockProgress, error)
	GetUserAnswers(ctx context.Context, in *dto.UserEventIds) (corr int64, total int64, err error)
	GetEventProgress(ctx context.Context, userId string, eventId string) (currTaskId string, currBlockId string, finished bool, currTaskStartTime time.Time, err error)
	GetBlockMaxPoints(ctx context.Context, blockId string) (int64, error)
	// PostAnswerBatch(ctx context.Context, in *dto.AnswerBatch) (*dto.MessageOut, error)
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
