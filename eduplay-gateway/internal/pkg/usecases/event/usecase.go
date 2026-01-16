package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	"log/slog"
)

type EventClient interface {
	SaveFile(ctx context.Context, in *dto.SaveFileIn) (*dto.MessageOut, error)
}

type UseCase struct {
	eventClient EventClient
	log         *slog.Logger
}

func New(log *slog.Logger, eventClient EventClient) *UseCase {
	return &UseCase{log: log, eventClient: eventClient}
}
