package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PostTask(ctx context.Context, req *eventModel.Task) (string, error) {
	const op = "event.UseCase.PostTask"

	s.log.With(slog.String("op", op)).Info("attempting to post task")

	corr, err := s.CheckTaskOptions(ctx, op, req)
	if !corr {
		return "", err
	}

	taskDto := eventModel.TaskToDto(req)

	ret, err := s.eventClient.PostTask(ctx, taskDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post task", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("task posted", slog.Any("event", ret))

	return ret.Message, nil
}
