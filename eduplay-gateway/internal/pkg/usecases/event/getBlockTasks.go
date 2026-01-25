package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetBlockTasks(ctx context.Context, req *eventModel.Id) (*eventModel.BlockTasksList, error) {
	const op = "event.UseCase.GetBlockTasks"

	s.log.With(slog.String("op", op)).Info("attempting to get block tasks")

	ret, err := s.eventClient.GetBlockTasks(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get block tasks", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got block tasks", slog.Any("event", ret))

	return eventModel.BlockTasksListFromDto(ret), nil
}
