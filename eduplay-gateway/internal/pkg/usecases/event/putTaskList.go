package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutTaskList(ctx context.Context, req *eventModel.PutTaskListIn) (string, error) {
	const op = "event.UseCase.PutTaskList"

	s.log.With(slog.String("op", op)).Info("attempting to update task list order")

	message, err := s.eventClient.PutTaskList(ctx, eventModel.PutTaskListInToDto(req))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to update task list order", slog.String("error", err.Error()))
		return "", err
	}

	return message.Message, nil
}
