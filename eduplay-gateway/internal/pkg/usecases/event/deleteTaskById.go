package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	"log/slog"
)

func (s *UseCase) DeleteTaskById(ctx context.Context, taskId string, eventId string) (string, error) {
	const op = "event.UseCase.DeleteTaskById"

	s.log.With(slog.String("op", op)).Info("attempting to delete task by id")

	ret, err := s.eventClient.DeleteTaskById(ctx, &eventDto.UserEventIds{
		EventId: eventId,
		UserId:  taskId,
	})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to delete task by id", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("deleted task", slog.Any("task", taskId))

	return ret.Message, nil
}
