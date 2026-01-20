package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	"log/slog"
)

type EventClient interface {
	SaveFile(ctx context.Context, in *dto.SaveFileIn) (*dto.MessageOut, error)
	PostEvent(ctx context.Context, in *dto.PostEventIn) (*dto.MessageOut, error)
	GetEvent(ctx context.Context, in *dto.Id) (*dto.PostEventIn, error)
	GetRole(ctx context.Context, in *dto.GetRoleIn) (*dto.GetRoleOut, error)
	GetGroups(ctx context.Context, in *dto.Id) (*dto.GetGroupsOut, error)
	GetCollaborators(ctx context.Context, in *dto.Id) (*dto.GetCollaboratorsOut, error)
	PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (*dto.MessageOut, error)
	GetEventBlocks(ctx context.Context, in *dto.Id) (*dto.GetEventBlocksOut, error)
}

type UseCase struct {
	eventClient EventClient
	log         *slog.Logger
}

func New(log *slog.Logger, eventClient EventClient) *UseCase {
	return &UseCase{log: log, eventClient: eventClient}
}
