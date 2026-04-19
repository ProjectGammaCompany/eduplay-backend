package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetUserAnswers(ctx context.Context, in *dto.UserEventIds) (*dto.UserAnswers, error) {
	const op = "Events.UseCase.GetUserAnswers"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting user answers")

	corr, total, err := a.storage.GetUserAnswers(ctx, in)
	if err != nil {
		log.Error("failed to get user answers in event", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	points, err := a.storage.GetUserStats(ctx, in.UserId, in.EventId)
	if err != nil {
		log.Error("failed to get user stats in event", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return &dto.UserAnswers{Correct: corr, Total: total, Points: points.Points}, nil
}
