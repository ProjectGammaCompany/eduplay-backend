package data_handler

import (
	"context"
	dto "eduplay-data/internal/generated"
	"fmt"

	"google.golang.org/grpc"
)

type UseCase interface {
	GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error)
	// SaveFile(ctx context.Context, in *dto.SaveFileIn) (string, error)
}

type Handler struct {
	dto.UnimplementedDataServer
	uc UseCase
}

func Register(gRPCServer *grpc.Server, uc UseCase) {
	dto.RegisterDataServer(gRPCServer, &Handler{uc: uc})
}

// func (h *Handler) SaveFile(ctx context.Context, in *dto.SaveFileIn) (*dto.MessageOut, error) {
// 	op := "SaveFile.Handler"

// 	message, err := h.uc.SaveFile(ctx, in)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return &dto.MessageOut{Message: message}, nil
// }

func (h *Handler) GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	op := "GetPublicEvents.Client"

	out, err := h.uc.GetPublicEvents(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}
