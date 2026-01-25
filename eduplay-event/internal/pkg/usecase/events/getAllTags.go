package event

import (
	"context"
	"log/slog"
	"strconv"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetAllTags(ctx context.Context) (*dto.Tags, error) {
	const op = "Events.UseCase.GetAllTags"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting all tags")

	tags, err := a.storage.GetAllTags(ctx)
	if err != nil {
		log.Error("failed to get all tags", err.Error(), slog.String("tags", strconv.Itoa(len(tags.Tags))))
		return nil, err
	}

	return tags, nil
}
