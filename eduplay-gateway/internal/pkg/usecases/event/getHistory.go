package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetHistory(ctx context.Context, filters *eventModel.EventBaseFilters) (*eventModel.GetPublicEventsOut, error) {
	const op = "event.UseCase.GetHistory"

	s.log.With(slog.String("op", op)).Info("attempting to get user history")

	events, err := s.eventClient.GetHistory(ctx, eventModel.EventBaseFiltersToDto(filters))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get user history", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got user history", slog.Any("event", len(events.Events)))

	return eventModel.GetPublicEventsOutFromDto(events), nil
}
