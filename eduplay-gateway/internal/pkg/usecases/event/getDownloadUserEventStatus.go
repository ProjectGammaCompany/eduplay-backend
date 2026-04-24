package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetDownloadUserEventStatus(ctx context.Context, req *eventModel.EventIds, userId string) (*eventModel.EventStatuses, error) {
	const op = "event.UseCase.GetDownloadUserEventStatus"

	s.log.With(slog.String("op", op)).Info("attempting to get user status for downloaded events")

	userEventStatuses := make([]eventModel.EventStatus, 0)

	for _, eventId := range req.EventIds {
		eventStatus := eventModel.EventStatus{}

		eventStatus.CompletedTasksInBlock = make([]string, 0)

		nextStage, err := s.eventClient.GetNextStage(ctx, &dto.UserEventIds{UserId: userId, EventId: eventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user status for downloaded events", slog.String("error", err.Error()))
			return nil, err
		}

		status, err := s.eventClient.GetUserStatus(ctx, &dto.UserEventIds{UserId: userId, EventId: eventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user status for downloaded events", slog.String("error", err.Error()))
			return nil, err
		}

		eventStatus.EventId = eventId
		eventStatus.Status = status.Status
		eventStatus.Type = nextStage.Type

		switch nextStage.Type {
		case "task":
			eventStatus.TaskId = nextStage.Task.TaskId
			eventStatus.BlockId = nextStage.Task.BlockId

			blockProgress, err := s.eventClient.GetBlockProgress(ctx, &dto.UserEventIds{UserId: userId, EventId: nextStage.Task.BlockId})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get user status for downloaded events", slog.String("error", err.Error()))
				return nil, err
			}

			eventStatus.PointsInBlock = blockProgress.PointsInBlock
			eventStatus.CompletedTasksInBlock = blockProgress.CompletedTasks
		case "block":
			eventStatus.BlockId = nextStage.Block.BlockId

			blockProgress, err := s.eventClient.GetBlockProgress(ctx, &dto.UserEventIds{UserId: userId, EventId: nextStage.Block.BlockId})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get user status for downloaded events", slog.String("error", err.Error()))
				return nil, err
			}

			eventStatus.PointsInBlock = blockProgress.PointsInBlock
			eventStatus.CompletedTasksInBlock = blockProgress.CompletedTasks
		}

		if status.Timestamp.AsTime().Format("02.01.2006 15:04:05.000") != "01.01.1970 00:00:00.000" {
			eventStatus.Timestamp = status.Timestamp.AsTime().Format("02.01.2006 15:04:05.000")
		}

		group, err := s.eventClient.GetUserGroup(ctx, &dto.UserEventIds{UserId: userId, EventId: eventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user status for downloaded events", slog.String("error", err.Error()))
			return nil, err
		}
		if group != nil && group.GroupId != "" {
			// TODO might change to group name
			eventStatus.GroupName = group.GroupId
		}

		event, err := s.eventClient.GetEvent(ctx, &dto.Id{Id: eventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user status for downloaded events", slog.String("error", err.Error()))
			return nil, err
		}

		eventStatus.LastEditionDate = event.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000")

		userEventStatuses = append(userEventStatuses, eventStatus)
	}
	s.log.With(slog.String("op", op)).Info("user status for downloaded events got", slog.Any("user", userId))

	return &eventModel.EventStatuses{EventStatuses: userEventStatuses}, nil
}
