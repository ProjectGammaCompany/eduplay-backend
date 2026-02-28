package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"fmt"
	"log/slog"
)

func (s *UseCase) PutTask(ctx context.Context, req *eventModel.Task) (*eventModel.PutTaskOut, error) {
	const op = "event.UseCase.PuttTask"

	s.log.With(slog.String("op", op)).Info("attempting to put task")

	taskDto := eventModel.TaskToDto(req)

	fmt.Println("Helloooooooo ehhhhhh?")

	ret, err := s.eventClient.PutTask(ctx, taskDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to put task", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("task put", slog.Any("event", ret))

	fmt.Println("Helloooooooo?")
	return eventModel.PutTaskOutFromDto(ret), nil
}
