package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetTaskById(ctx context.Context, req *eventModel.Id) (*eventModel.Task, error) {
	const op = "event.UseCase.GetTaskById"

	s.log.With(slog.String("op", op)).Info("attempting to get task by id")

	ret, err := s.eventClient.GetTaskById(ctx, &eventDto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get task by id", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got tasks", slog.Any("task", req.Id))

	return eventModel.TaskFromDto(ret), nil
}
