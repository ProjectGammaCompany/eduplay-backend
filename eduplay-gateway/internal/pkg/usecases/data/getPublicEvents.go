package data

import (
	"context"
	dataModel "eduplay-gateway/internal/lib/models/data"
	"log/slog"
)

func (s *UseCase) GetPublicEvents(ctx context.Context, filters *dataModel.EventBaseFilters) (*dataModel.GetPublicEventsOut, error) {
	const op = "data.UseCase.GetPublicEvents"

	s.log.With(slog.String("op", op)).Info("attempting to get public events")

	events, err := s.dataClient.GetPublicEvents(ctx, dataModel.EventBaseFiltersToDto(filters))
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get public events", slog.String("error", err.Error()))
		return nil, err
	}

	s.log.With(slog.String("op", op)).Info("got public events", slog.Any("event", len(events.Events)))

	return dataModel.GetPublicEventsOutFromDto(events), nil
}
