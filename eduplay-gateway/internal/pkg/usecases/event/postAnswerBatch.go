package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PostAnswerBatch(ctx context.Context, req *eventModel.AnswerBatch) (string, error) {
	const op = "event.UseCase.PostAnswerBatch"

	s.log.With(slog.String("op", op)).Info("attempting to post answer batch")

	answerBatchDto, err := eventModel.AnswerBatchToDto(req)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to convert answer batch to dto", slog.String("error", err.Error()))
		return "", err
	}

	ret, err := s.eventClient.PostAnswerBatch(ctx, answerBatchDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post answer batch", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("answer batch posted", slog.Any("user", req.UserId), slog.Any("event", req.EventId))

	return ret.Message, nil
}
