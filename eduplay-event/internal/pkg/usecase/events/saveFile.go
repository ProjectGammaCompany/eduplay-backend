package event

import (
	"context"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) SaveFile(ctx context.Context, in *dto.SaveFileIn) (string, error) {
	const op = "Events.UseCase.SaveFile"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("saving file")

	message, err := a.storage.SaveFile(ctx, in.Filename, in.FileUUID)
	if err != nil {
		log.Error("failed to save file info", err.Error(), slog.String("filename", in.Filename))
		return "", err
	}

	return message, nil
}
