package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetOwnedEvents(ctx context.Context, filters *eventModel.EventBaseFilters) (*eventModel.GetPublicEventsOut, error) {
	const op = "event.UseCase.GetOwnedEvents"

	s.log.With(slog.String("op", op)).Info("attempting to get user owned events")

	events, err := s.eventClient.GetOwnedEvents(ctx, eventModel.EventBaseFiltersToDto(filters))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get user owned events", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got user owned events", slog.Any("event", len(events.Events)))

	return eventModel.GetPublicEventsOutFromDto(events), nil
}
