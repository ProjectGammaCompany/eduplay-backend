package event

import (
	"context"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetAllTags(ctx context.Context) (eventModel.Tags, error) {
	const op = "event.UseCase.GetAllTags"

	s.log.With(slog.String("op", op)).Info("attempting to get all tags")

	tags, err := s.eventClient.GetAllTags(ctx)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get all tags", slog.String("error", err.Error()))
		return eventModel.Tags{}, err
	}

	ret := eventModel.TagsFromDto(tags.Tags)

	s.log.With(slog.String("op", op)).Info("got all tags", slog.Any("tags", len(ret.Tags)))

	return ret, nil
}
