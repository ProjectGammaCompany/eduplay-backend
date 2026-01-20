package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetUserFavorites(ctx context.Context, filters *eventModel.EventBaseFilters) (*eventModel.GetPublicEventsOut, error) {
	const op = "event.UseCase.GetUserFavorites"

	s.log.With(slog.String("op", op)).Info("attempting to get user favorites")

	events, err := s.eventClient.GetUserFavorites(ctx, eventModel.EventBaseFiltersToDto(filters))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get user favorites", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got use favorites", slog.Any("event", len(events.Events)))

	return eventModel.GetPublicEventsOutFromDto(events), nil
}
