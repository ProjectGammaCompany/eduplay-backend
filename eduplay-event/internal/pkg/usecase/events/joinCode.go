package event

import (
	"context"
	"crypto/rand"
	dto "eduplay-event/internal/generated"
	errs "eduplay-event/internal/storage"
	"encoding/base32"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var encoding = base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ23456789").WithPadding(base32.NoPadding)

func GenerateJoinCode(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return encoding.EncodeToString(b)[:length], nil
}

func (a *UseCase) PostJoinCode(ctx context.Context, in *dto.Id) (*dto.JoinCode, error) {
	const op = "Events.UseCase.PostJoinCode"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to save event join code")

	for i := 0; i < 5; i++ {
		joinCode, err := GenerateJoinCode(6)
		if err != nil {
			log.Error("failed to generate join code", slog.String("error", err.Error()))
			return nil, err
		}

		time, err := a.storage.InsertJoinCode(ctx, in.Id, joinCode)
		if err != nil {
			if err == errs.ErrJoinCodeNotUnique {
				continue
			}
			log.Error("failed to save join code", slog.String("error", err.Error()))
			return nil, err
		}

		return &dto.JoinCode{EventId: in.Id, JoinCode: joinCode, ExpiresAt: timestamppb.New(*time)}, nil
	}
	return nil, errs.ErrJoinCodeRetryFailed
}

func (a *UseCase) GetJoinCode(ctx context.Context, in *dto.Id) (*dto.JoinCode, error) {
	const op = "Events.UseCase.GetJoinCode"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get event join code")

	joinCode, err := a.storage.GetJoinCode(ctx, in.Id)
	if err != nil {
		if err == errs.ErrJoinCodeExpired {
			joinCode, err := a.PostJoinCode(ctx, in)
			if err != nil {
				return nil, err
			}
			return joinCode, nil
		}
		log.Error("failed to get join code", err.Error(), slog.String("event", in.Id))
		return nil, err
	}

	return joinCode, nil
}
