package data

import (
	"context"
	"log/slog"

	// "eduplay-event/internal/model"
	// dto "eduplay-data/internal/generated"
	eventDto "eduplay-data/internal/generated/clients/event"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

type storage interface {
	// SaveFile(ctx context.Context, fileName string, fileKey string, fileUUID string) (string, error)
}

type EventClient interface {
	SaveFile(ctx context.Context, in *eventDto.SaveFileIn) (*eventDto.MessageOut, error)
	GetPublicEvents(ctx context.Context, in *eventDto.EventBaseFilters) (*eventDto.GetPublicEventsOut, error)
}

type UseCase struct {
	log      *slog.Logger
	storage  storage
	evClient EventClient
	secret   string
}

func New(
	log *slog.Logger,
	st storage,
	evClient EventClient,
	secret string,
) *UseCase {
	return &UseCase{
		log:      log,
		storage:  st,
		evClient: evClient,
		secret:   secret,
	}
}
