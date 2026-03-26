package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostComplaint(ctx context.Context, in *dto.PostComplaintIn) (string, error) {
	const op = "Events.UseCase.PostComplaint"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("posting a complaint")

	ret, err := a.storage.PostComplaint(ctx, in)
	if err != nil {
		log.Error("failed to post complaint", err.Error(), slog.String("complaint reason", in.Reason))
		return "", err
	}

	return ret, nil
}
