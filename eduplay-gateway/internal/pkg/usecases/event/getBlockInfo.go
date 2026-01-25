package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetBlockInfo(ctx context.Context, req *eventModel.Id) (*eventModel.PostEventBlockIn, error) {
	const op = "event.UseCase.GetBlockInfo"

	s.log.With(slog.String("op", op)).Info("attempting to ge block info")

	ret, err := s.eventClient.GetBlockInfo(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get block info", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got block info", slog.Any("event", ret))

	return eventModel.PostEventBlockFromDto(ret), nil
}
