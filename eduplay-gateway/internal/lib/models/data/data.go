package eventModel

import (
	dto "eduplay-gateway/internal/generated/clients/data"
	// 	"time"
	// 	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventBaseFilters struct {
	Page            int64    `json:"page"`
	MaxOnPage       int64    `json:"maxOnPage"`
	Tags            []string `json:"tags"`
	DecliningRating bool     `json:"decliningRating"`
	Territorialized bool     `json:"territorialized"`
	Active          bool     `json:"active"`
	Favorites       bool     `json:"favorites"`
	UserId          string   `json:"userId"`
	Title           string   `json:"title"`
}

func EventBaseFiltersToDto(in *EventBaseFilters) *dto.EventBaseFilters {
	return &dto.EventBaseFilters{
		Page:            in.Page,
		MaxOnPage:       in.MaxOnPage,
		Tags:            in.Tags,
		DecliningRating: in.DecliningRating,
		Territorialized: in.Territorialized,
		Active:          in.Active,
		Favorites:       in.Favorites,
		UserId:          in.UserId,
		Title:           in.Title,
	}
}

type GetPublicEvent struct {
	EventId         string `json:"id"`
	Title           string `json:"title" validate:"required"`
	Description     string `json:"description"`
	Tags            []Tag  `json:"tags"`
	Cover           string `json:"cover"`
	LastEditionDate string `json:"lastEditionDate"`
	Rate            int64  `json:"rate"`
	Favorite        bool   `json:"favorite"`
}

type GetPublicEventsOut struct {
	Events []*GetPublicEvent `json:"events"`
}

func GetPublicEventFromDto(in *dto.GetPublicEvent) *GetPublicEvent {
	event := &GetPublicEvent{
		EventId:         in.EventId,
		Title:           in.Title,
		Description:     in.Description,
		Tags:            TagsFromDto(in.Tags).Tags,
		Cover:           in.Cover,
		LastEditionDate: in.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000"),
		Rate:            in.Rate,
		Favorite:        in.Favorite,
	}
	if event.LastEditionDate == "01.01.1970 00:00:00.000" {
		event.LastEditionDate = ""
	}

	return event
}

func GetPublicEventsOutFromDto(in *dto.GetPublicEventsOut) *GetPublicEventsOut {
	events := make([]*GetPublicEvent, len(in.Events))
	for i, event := range in.Events {
		events[i] = GetPublicEventFromDto(event)
	}

	return &GetPublicEventsOut{
		Events: events,
	}
}

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
}

func TagsFromDto(in []*dto.Tag) Tags {
	tags := make([]Tag, len(in))
	for i, tag := range in {
		tags[i] = Tag{
			Id:   tag.Id,
			Name: tag.Name,
		}
	}

	return Tags{
		Tags: tags,
	}
}
