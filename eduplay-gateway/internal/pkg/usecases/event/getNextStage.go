package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	errs "eduplay-gateway/internal/storage"
	"log/slog"
)

func (s *UseCase) GetNextStage(ctx context.Context, in *eventModel.UserEventIds) (*eventModel.NextStageInfo, error) {
	const op = "event.UseCase.GetNextStage"

	s.log.With(slog.String("op", op)).Info("attempting to get next stage")

	ret, err := s.eventClient.GetNextStage(ctx, eventModel.UserEventIdsToDto(in))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get next stage", slog.String("error", err.Error()))
		return nil, err
	}

	if ret.Type == "not found" {
		return nil, errs.ErrNotFound
	}

	s.log.With(slog.String("op", op)).Info("event get next stage", slog.Any("event", ret))

	return eventModel.NextStageInfoFromDto(ret), nil
}
