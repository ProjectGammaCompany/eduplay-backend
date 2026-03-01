package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	"log/slog"
)

func (a *UseCase) SaveFile(ctx context.Context, fileName string, fileKey string, fileUUID string) (string, error) {
	const op = "Events.UseCase.SaveFile"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to save file")

	ret, err := a.eventClient.SaveFile(ctx, &dto.SaveFileIn{Filename: fileName, FileKey: fileKey, FileUUID: fileUUID})
	if err != nil {
		log.Error("failed to save file", err.Error(), slog.String("filename", fileName))
		return "", err
	}

	return ret.Message, nil
}
