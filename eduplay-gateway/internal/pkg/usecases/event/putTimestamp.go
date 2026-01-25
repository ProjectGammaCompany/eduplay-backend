package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutTimestamp(ctx context.Context, in *eventModel.PutTimestampIn) (string, error) {
	const op = "event.UseCase.PutTimestamp"

	s.log.With(slog.String("op", op)).Info("attempting to put timestamp")

	putTimestampDto, err := eventModel.PutTimestampInToDto(in)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to convert event to dto", slog.String("error", err.Error()))
		return "", err
	}

	ret, err := s.eventClient.PutTimestamp(ctx, putTimestampDto)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to put timestamp", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("event put timestamp", slog.Any("event", ret))

	return ret.Message, nil
}
