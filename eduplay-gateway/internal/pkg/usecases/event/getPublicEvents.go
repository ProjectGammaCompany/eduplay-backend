package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetPublicEvents(ctx context.Context, filters *eventModel.EventBaseFilters) (*eventModel.GetPublicEventsOut, error) {
	const op = "event.UseCase.GetPublicEvents"

	s.log.With(slog.String("op", op)).Info("attempting to get public events")

	events, err := s.eventClient.GetPublicEvents(ctx, eventModel.EventBaseFiltersToDto(filters))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get public events", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got public events", slog.Any("event", len(events.Events)))

	return eventModel.GetPublicEventsOutFromDto(events), nil
}
