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
	GetEvent(ctx context.Context, in *dto.Id) (*dto.PostEventIn, error)
	GetRole(ctx context.Context, in *dto.GetRoleIn) (*dto.GetRoleOut, error)
	GetGroups(ctx context.Context, in *dto.Id) (*dto.GetGroupsOut, error)
	GetCollaborators(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error)
	PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	GetEventBlocks(ctx context.Context, in *dto.Id) (*dto.GetEventBlocksOut, error)
	GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
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
