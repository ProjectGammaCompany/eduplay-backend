package event

import (
	"context"
	dto "eduplay-event/internal/generated"
	"log/slog"
)

func (a *UseCase) GetEditorUserStatsTask(ctx context.Context, in *dto.UserEventIds) (*dto.EditorStatsTask, error) {
	const op = "Events.UseCase.GetEditorUserStatsTask"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting editor user stats task")

	ret, err := a.storage.GetEditorUserStatsTask(ctx, in.UserId, in.EventId)
	if err != nil {
		log.Error("failed to get editor user stats task", err.Error(), slog.String("task", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return ret, nil
}
