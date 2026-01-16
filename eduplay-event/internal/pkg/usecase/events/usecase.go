package event

import (
	"context"
	"log/slog"
	// "eduplay-event/internal/model"
)

type storage interface {
	SaveFile(ctx context.Context, fileName string, fileUUID string) (string, error)
}

type UseCase struct {
	log     *slog.Logger
	storage storage
	secret  string
}

func New(
	log *slog.Logger,
	st storage,
	secret string,
) *UseCase {
	return &UseCase{
		log:     log,
		storage: st,
		secret:  secret,
	}
}
