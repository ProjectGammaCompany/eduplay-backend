package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) PutFavorite(ctx context.Context, req *eventModel.PutFavorite) (string, error) {
	const op = "event.UseCase.PutFavorite"

	s.log.With(slog.String("op", op)).Info("attempting to put favorite")

	favorite := &dto.PutFavoriteIn{
		UserId:   req.UserId,
		EventId:  req.EventId,
		Favorite: req.Favorite,
	}

	ret, err := s.eventClient.PutFavorite(ctx, favorite)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to put favorite", slog.String("error", err.Error()))
		return "", err
	}

	s.log.With(slog.String("op", op)).Info("event put favorite", slog.Any("event", ret))

	return ret.Message, nil
}
