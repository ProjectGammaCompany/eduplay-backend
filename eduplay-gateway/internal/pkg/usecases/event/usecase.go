package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	userDto "eduplay-gateway/internal/generated/clients/user"
	eventModel "eduplay-gateway/internal/lib/models/event"
	errs "eduplay-gateway/internal/storage"
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
	GetUserStatus(ctx context.Context, in *dto.UserEventIds) (*dto.UserStatus, error)
	GetGroupUsers(ctx context.Context, in *dto.Id) (*dto.GetGroupUsersOut, error)
	GetUserStats(ctx context.Context, in *dto.UserEventIds) (*dto.User, error)
	GetUserGroup(ctx context.Context, in *dto.UserEventIds) (*dto.GetUserGroupOut, error)
	GetEventUsers(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error)
	PostComplaint(ctx context.Context, in *dto.PostComplaintIn) (*dto.MessageOut, error)
	GetJoinCode(ctx context.Context, in *dto.Id) (*dto.JoinCode, error)
	GetEventByJoinCode(ctx context.Context, in *dto.Id) (*dto.Id, error)
	GetEventUserRating(ctx context.Context, in *dto.UserEventIds) (*dto.MessageOut, error)
	PostParticipant(ctx context.Context, in *dto.PostParticipantIn) (*dto.MessageOut, error)
	PostRate(ctx context.Context, in *dto.Rate) (*dto.MessageOut, error)
	GetBlockProgress(ctx context.Context, in *dto.UserEventIds) (*dto.BlockProgress, error)
}

type UserClient interface {
	GetProfile(ctx context.Context, userId string) (*userDto.Profile, error)
	GetProfileByLogin(ctx context.Context, login string) (*userDto.Profile, error)
}

type UseCase struct {
	eventClient EventClient
	userClient  UserClient
	log         *slog.Logger
}

func New(log *slog.Logger, eventClient EventClient, userClient UserClient) *UseCase {
	return &UseCase{log: log, eventClient: eventClient, userClient: userClient}
}

func (s *UseCase) CheckTaskOptions(ctx context.Context, op string, req *eventModel.Task) (bool, error) {
	switch req.TaskType {
	case 0:
		if len(req.Options) > 0 {
			s.log.With(slog.String("op", op)).Error("task type does not support options",
				slog.String("error", "task type does not support options"))
			return false, errs.ErrInfoSegmentAnswerIncorrect
		}
	case 1:
		count := 0
		for _, option := range req.Options {
			if option.IsCorrect {
				count++
			}
		}
		if count != 1 {
			s.log.With(slog.String("op", op)).Error("single choice tasks do not support such options",
				slog.String("error", "single choice tasks do not support such options"))
			return false, errs.ErrSingleChoiceAnswerIncorrect
		}
	case 2:
		count := 0
		for _, option := range req.Options {
			if option.IsCorrect {
				count++
			}
		}
		if count < 1 {
			s.log.With(slog.String("op", op)).Error("multiple choice tasks do not support such options",
				slog.String("error", "multiple choice tasks do not support such options"))
			return false, errs.ErrMultipleChoiceAnswerIncorrect
		}
	case 3:
		if len(req.Options) != 1 {
			s.log.With(slog.String("op", op)).Error("text tasks do not support such options",
				slog.String("error", "text tasks do not support such options"))
			return false, errs.ErrTextAnswerIncorrect
		}
	case 4:
		if len(req.Options) != 1 {
			s.log.With(slog.String("op", op)).Error("text tasks do not support such options",
				slog.String("error", "text tasks do not support such options"))
			return false, errs.ErrTextAnswerIncorrect
		}
	}
	return true, nil
}
