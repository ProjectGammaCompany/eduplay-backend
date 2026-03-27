package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	errs "eduplay-gateway/internal/storage"
	"log/slog"
)

func (s *UseCase) PostComplaint(ctx context.Context, req *eventModel.Complaint) (string, error) {
	const op = "event.UseCase.PostComplaint"

	s.log.With(slog.String("op", op)).Info("attempting to post complaint")

	complaintDto := &eventDto.PostComplaintIn{
		Reason:  req.Reason,
		EventId: req.EventId,
		UserId:  req.UserId,
	}

	ret, err := s.eventClient.PostComplaint(ctx, complaintDto)

	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to post complaint", slog.String("error", err.Error()))
		return "", err
	}

	if ret == nil || ret.Message == "" {
		s.log.With(slog.String("op", op)).Error("failed to post complaint", slog.String("error", errs.ErrNoRows.Error()))
		return "", errs.ErrNoRows
	}

	s.log.With(slog.String("op", op)).Info("complaint posted", slog.Any("event", ret))

	return ret.Message, nil
}
