package event

import (
	"context"
	"log/slog"

	// "eduplay-event/internal/model"
	dto "eduplay-event/internal/generated"
)

type storage interface {
	SaveFile(ctx context.Context, fileName string, fileUUID string) (string, error)
	PostEvent(ctx context.Context, in *dto.PostEventIn) (string, error)
	GetEvent(ctx context.Context, id string) (*dto.PostEventIn, error)
	GetRole(ctx context.Context, userId string, eventId string) (int64, error)
	GetGroups(ctx context.Context, eventId string) (*dto.GetGroupsOut, error)
	GetCollaborators(ctx context.Context, eventId string) (*dto.GetCollaboratorsOut, error)
	PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error)
	GetEventBlocks(ctx context.Context, eventId string) (*dto.GetEventBlocksOut, error)
	GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
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
