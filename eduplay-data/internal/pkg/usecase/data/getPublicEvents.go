package data

import (
	"context"
	"log/slog"

	dto "eduplay-data/internal/generated"
	eventDto "eduplay-data/internal/generated/clients/event"
)

func (a *UseCase) GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "Events.UseCase.GetPublicEvents"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting public events")

	filters := &eventDto.EventBaseFilters{
		Page:            in.Page,
		MaxOnPage:       in.MaxOnPage,
		UserId:          in.UserId,
		DecliningRating: in.DecliningRating,
		Territorialized: in.Territorialized,
		Active:          in.Active,
		Favorites:       in.Favorites,
		Tags:            in.Tags,
		Title:           in.Title,
	}

	events, err := a.evClient.GetPublicEvents(ctx, filters)
	if err != nil {
		log.Error("failed to get public events", err.Error(), slog.String("event", in.String()))
		return nil, err
	}

	allEvents := &dto.GetPublicEventsOut{}
	dataEvents := make([]*dto.GetPublicEvent, 0)
	for _, event := range events.Events {
		dataTags := make([]*dto.Tag, 0)
		for _, tag := range event.Tags {
			dataTags = append(dataTags, &dto.Tag{
				Id:   tag.Id,
				Name: tag.Name,
			})
		}
		dataEvent := &dto.GetPublicEvent{
			EventId:         event.EventId,
			Title:           event.Title,
			Description:     event.Description,
			Cover:           event.Cover,
			LastEditionDate: event.LastEditionDate,
			Tags:            dataTags,
			Rate:            event.Rate,
			Favorite:        event.Favorite,
		}
		dataEvents = append(dataEvents, dataEvent)
	}

	allEvents.Events = dataEvents

	return allEvents, nil
}
