package sign_up_user

import (
	"context"

	"fmt"

	dto "eduplay-event/internal/generated"

	"google.golang.org/grpc"
)

type UseCase interface {
	SaveFile(ctx context.Context, in *dto.SaveFileIn) (string, error)
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
