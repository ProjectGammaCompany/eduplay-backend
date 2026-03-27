package sign_up_user

import (
	"context"

	"fmt"

	dto "eduplay-event/internal/generated"

	"google.golang.org/grpc"
)

type UseCase interface {
	SaveFile(ctx context.Context, in *dto.SaveFileIn) (string, error)
	PostEvent(ctx context.Context, in *dto.PostEventIn) (string, error)
	PutEvent(ctx context.Context, in *dto.PutEventIn) (*dto.GetGroupsOut, error)
	GetEvent(ctx context.Context, in *dto.Id) (*dto.PostEventIn, error)
	GetRole(ctx context.Context, in *dto.GetRoleIn) (*dto.GetRoleOut, error)
	GetGroups(ctx context.Context, in *dto.Id) (*dto.GetGroupsOut, error)
	PutGroups(ctx context.Context, in *dto.PutListIn) (string, error)
	PutTaskList(ctx context.Context, in *dto.PutListIn) (string, error)
	PutBlockList(ctx context.Context, in *dto.PutListIn) (string, error)
	GetCollaborators(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error)
	PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	PutEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	PutEventBlockName(ctx context.Context, in *dto.Tag) (string, error)
	GetEventBlocks(ctx context.Context, in *dto.Id) (*dto.GetEventBlocksOut, error)
	GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetOwnedEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetHistory(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	PutFavorite(ctx context.Context, in *dto.PutFavoriteIn) (string, error)
	GetAllTags(ctx context.Context) (*dto.Tags, error)
	PostTask(ctx context.Context, in *dto.Task) (string, error)
	PutTask(ctx context.Context, in *dto.Task) (*dto.PutTaskOut, error)
	PostBlockCondition(ctx context.Context, in *dto.Condition) (*dto.PostConditionOut, error)
	PutBlockCondition(ctx context.Context, in *dto.Condition) (*dto.MessageOut, error)
	DeleteBlockCondition(ctx context.Context, in *dto.Id) (string, error)
	GetBlockInfo(ctx context.Context, in *dto.Id) (*dto.PostEventBlockIn, error)
	GetBlockConditions(ctx context.Context, in *dto.Id) (*dto.BlockInfo, error)
	GetBlockTasks(ctx context.Context, in *dto.Id) (*dto.Tasks, error)
	GetTaskById(ctx context.Context, in *dto.Id) (*dto.Task, error)
	DeleteTaskById(ctx context.Context, in *dto.Id) (string, error)
	PostAnswer(ctx context.Context, in *dto.Answer) (*dto.Answer, error)
	DeleteBlockById(ctx context.Context, in *dto.Id) (string, error)
	DeleteEventById(ctx context.Context, in *dto.Id) (string, error)
	GetPublicEvent(ctx context.Context, in *dto.UserEventIds) (*dto.GetPublicEvent, error)
	PutNextStage(ctx context.Context, stage *dto.EventBlockTaskUserIds) (string, error)
	GetNextStage(ctx context.Context, in *dto.UserEventIds) (*dto.NextStageInfo, error)
	PutTimestamp(ctx context.Context, in *dto.PutTimestampIn) (string, error)
	GetUserStatus(ctx context.Context, in *dto.UserEventIds) (*dto.MessageOut, error)
	GetUserStats(ctx context.Context, in *dto.UserEventIds) (*dto.User, error)
	GetGroupUsers(ctx context.Context, in *dto.Id) (*dto.GetGroupUsersOut, error)
	GetUserGroup(ctx context.Context, in *dto.UserEventIds) (*dto.GetUserGroupOut, error)
	GetEventUsers(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error)
	PostComplaint(ctx context.Context, in *dto.PostComplaintIn) (string, error)
	GetJoinCode(ctx context.Context, in *dto.Id) (*dto.JoinCode, error)
}

type Handler struct {
	dto.UnimplementedEventsServer
	uc UseCase
}

func Register(gRPCServer *grpc.Server, uc UseCase) {
	dto.RegisterEventsServer(gRPCServer, &Handler{uc: uc})
}

func (h *Handler) SaveFile(ctx context.Context, in *dto.SaveFileIn) (*dto.MessageOut, error) {
	op := "SaveFile.Handler"

	message, err := h.uc.SaveFile(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) PostEvent(ctx context.Context, in *dto.PostEventIn) (*dto.MessageOut, error) {
	op := "PostEvent.Handler"

	id, err := h.uc.PostEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: id}, nil
}

func (h *Handler) PutEvent(ctx context.Context, in *dto.PutEventIn) (*dto.GetGroupsOut, error) {
	op := "PutEvent.Handler"

	out, err := h.uc.PutEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (h *Handler) GetEvent(ctx context.Context, in *dto.Id) (*dto.PostEventIn, error) {
	op := "GetEvent.Handler"

	event, err := h.uc.GetEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return event, nil
}

func (h *Handler) GetRole(ctx context.Context, in *dto.GetRoleIn) (*dto.GetRoleOut, error) {
	op := "GetRole.Handler"

	role, err := h.uc.GetRole(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (h *Handler) GetGroups(ctx context.Context, in *dto.Id) (*dto.GetGroupsOut, error) {
	op := "GetGroups.Handler"

	groups, err := h.uc.GetGroups(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return groups, nil
}

func (h *Handler) PutGroups(ctx context.Context, in *dto.PutListIn) (*dto.MessageOut, error) {
	op := "PutGroups.Handler"

	message, err := h.uc.PutGroups(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) PutTaskList(ctx context.Context, in *dto.PutListIn) (*dto.MessageOut, error) {
	op := "PutTaskList.Handler"

	message, err := h.uc.PutTaskList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) PutBlockList(ctx context.Context, in *dto.PutListIn) (*dto.MessageOut, error) {
	op := "PutBlockList.Handler"

	message, err := h.uc.PutBlockList(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) GetCollaborators(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error) {
	op := "GetCollaborators.Handler"

	collaborators, err := h.uc.GetCollaborators(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return collaborators, nil
}

func (h *Handler) PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (*dto.MessageOut, error) {
	op := "PostEventBlock.Handler"

	id, err := h.uc.PostEventBlock(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: id}, nil
}

func (h *Handler) PutEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (*dto.MessageOut, error) {
	op := "PutEventBlock.Handler"

	id, err := h.uc.PutEventBlock(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: id}, nil
}

func (h *Handler) PutEventBlockName(ctx context.Context, in *dto.Tag) (*dto.MessageOut, error) {
	op := "PutEventBlock.Handler"

	message, err := h.uc.PutEventBlockName(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) GetEventBlocks(ctx context.Context, in *dto.Id) (*dto.GetEventBlocksOut, error) {
	op := "GetEventBlocks.Handler"

	blocks, err := h.uc.GetEventBlocks(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return blocks, nil
}

func (h *Handler) GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	op := "GetPublicEvents.Handler"

	events, err := h.uc.GetPublicEvents(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (h *Handler) GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	op := "GetUserFavorites.Handler"

	events, err := h.uc.GetUserFavorites(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (h *Handler) GetOwnedEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	op := "GetOwnedEvents.Handler"

	events, err := h.uc.GetOwnedEvents(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (h *Handler) GetHistory(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	op := "GetHistory.Handler"

	events, err := h.uc.GetHistory(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (h *Handler) PutFavorite(ctx context.Context, in *dto.PutFavoriteIn) (*dto.MessageOut, error) {
	op := "PutFavorite.Handler"

	message, err := h.uc.PutFavorite(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) GetAllTags(ctx context.Context, in *dto.Empty) (*dto.Tags, error) {
	op := "GetAllTags.Handler"

	tags, err := h.uc.GetAllTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tags, nil
}

func (h *Handler) PostTask(ctx context.Context, in *dto.Task) (*dto.MessageOut, error) {
	op := "PostTask.Handler"

	id, err := h.uc.PostTask(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: id}, nil
}

func (h *Handler) PutTask(ctx context.Context, in *dto.Task) (*dto.PutTaskOut, error) {
	op := "PutTask.Handler"

	ret, err := h.uc.PutTask(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) PostBlockCondition(ctx context.Context, in *dto.Condition) (*dto.PostConditionOut, error) {
	op := "PostBlockCondition.Handler"

	ret, err := h.uc.PostBlockCondition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) PutBlockCondition(ctx context.Context, in *dto.Condition) (*dto.MessageOut, error) {
	op := "PutBlockCondition.Handler"

	message, err := h.uc.PutBlockCondition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

func (h *Handler) GetBlockInfo(ctx context.Context, in *dto.Id) (*dto.PostEventBlockIn, error) {
	op := "GetBlockInfo.Handler"

	ret, err := h.uc.GetBlockInfo(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) GetBlockConditions(ctx context.Context, in *dto.Id) (*dto.BlockInfo, error) {
	op := "GetBlockConditions.Handler"

	ret, err := h.uc.GetBlockConditions(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) GetBlockTasks(ctx context.Context, in *dto.Id) (*dto.Tasks, error) {
	op := "GetBlockTasks.Handler"

	ret, err := h.uc.GetBlockTasks(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) GetTaskById(ctx context.Context, in *dto.Id) (*dto.Task, error) {
	op := "GetTaskById.Handler"

	ret, err := h.uc.GetTaskById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) DeleteTask(ctx context.Context, in *dto.Id) (*dto.MessageOut, error) {
	op := "DeleteTaskById.Handler"

	message, err := h.uc.DeleteTaskById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) PostAnswer(ctx context.Context, in *dto.Answer) (*dto.Answer, error) {
	op := "PostAnswer.Handler"

	ret, err := h.uc.PostAnswer(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) DeleteBlockById(ctx context.Context, in *dto.Id) (*dto.MessageOut, error) {
	op := "DeleteBlockById.Handler"

	message, err := h.uc.DeleteBlockById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) DeleteEventById(ctx context.Context, in *dto.Id) (*dto.MessageOut, error) {
	op := "DeleteEventById.Handler"

	message, err := h.uc.DeleteEventById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) DeleteBlockCondition(ctx context.Context, in *dto.Id) (*dto.MessageOut, error) {
	op := "DeleteBlockCondition.Handler"

	message, err := h.uc.DeleteBlockCondition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) GetEventForUser(ctx context.Context, in *dto.UserEventIds) (*dto.GetPublicEvent, error) {
	op := "GetPublicEvent.Handler"

	ret, err := h.uc.GetPublicEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) PutNextStage(ctx context.Context, in *dto.EventBlockTaskUserIds) (*dto.MessageOut, error) {
	op := "PutNextStage.Handler"

	message, err := h.uc.PutNextStage(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) GetNextStage(ctx context.Context, in *dto.UserEventIds) (*dto.NextStageInfo, error) {
	op := "GetNextStage.Handler"

	ret, err := h.uc.GetNextStage(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) PutTimestamp(ctx context.Context, in *dto.PutTimestampIn) (*dto.MessageOut, error) {
	op := "PutTimestamp.Handler"

	message, err := h.uc.PutTimestamp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) GetUserStatus(ctx context.Context, in *dto.UserEventIds) (*dto.MessageOut, error) {
	op := "GetUserStatus.Handler"

	ret, err := h.uc.GetUserStatus(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ret, nil
}

func (h *Handler) GetUserStats(ctx context.Context, in *dto.UserEventIds) (*dto.User, error) {
	op := "GetUserStats.Handler"
	out, err := h.uc.GetUserStats(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (h *Handler) GetGroupUsers(ctx context.Context, in *dto.Id) (*dto.GetGroupUsersOut, error) {
	op := "GetGroupUsers.Handler"
	out, err := h.uc.GetGroupUsers(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (h *Handler) GetUserGroup(ctx context.Context, in *dto.UserEventIds) (*dto.GetUserGroupOut, error) {
	op := "GetUserGroup.Handler"
	out, err := h.uc.GetUserGroup(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (h *Handler) GetEventUsers(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error) {
	op := "GetEventUsers.Handler"
	out, err := h.uc.GetEventUsers(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (h *Handler) PostComplaint(ctx context.Context, in *dto.PostComplaintIn) (*dto.MessageOut, error) {
	op := "PostComplaint.Handler"

	fmt.Println(op)

	message, err := h.uc.PostComplaint(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.MessageOut{Message: message}, nil
}

func (h *Handler) GetJoinCode(ctx context.Context, in *dto.Id) (*dto.JoinCode, error) {
	op := "GetJoinCode.Handler"
	out, err := h.uc.GetJoinCode(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}
