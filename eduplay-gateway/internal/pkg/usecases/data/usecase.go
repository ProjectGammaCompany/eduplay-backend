package data

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/data"
	"log/slog"
)

type DataClient interface {
	GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
}

type UseCase struct {
	dataClient DataClient
	log        *slog.Logger
}

func New(log *slog.Logger, dataClient DataClient) *UseCase {
	return &UseCase{log: log, dataClient: dataClient}
}
