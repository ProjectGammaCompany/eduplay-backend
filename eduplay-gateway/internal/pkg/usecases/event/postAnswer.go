package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PostAnswer(ctx context.Context, req *eventModel.Answer) (*eventModel.Answer, error) {
	const op = "event.UseCase.PostAnswer"

	s.log.With(slog.String("op", op)).Info("attempting to post answer")

	ret, err := s.eventClient.PostAnswer(ctx, eventModel.AnswerToDto(req))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post answer", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("posted answer", slog.Any("task", ret.TaskId))

	return eventModel.AnswerFromDto(ret), nil
}
