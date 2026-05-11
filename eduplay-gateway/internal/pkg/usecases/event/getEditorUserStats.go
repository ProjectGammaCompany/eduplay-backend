package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetEditorUserStats(ctx context.Context, req *eventModel.UserEventIds) (*eventModel.EditorUserStats, error) {
	const op = "event.UseCase.GetEditorUserStats"

	s.log.With(slog.String("op", op)).Info("attempting to get editor stats")

	editorStats := eventModel.EditorUserStats{}

	blocks, err := s.eventClient.GetEventBlocks(ctx, &dto.Id{Id: req.EventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return nil, err
	}

	blockModels := make([]eventModel.TaskBlock, 0)

	for _, block := range blocks.Blocks {
		tasks, err := s.eventClient.GetBlockTasks(ctx, &dto.Id{Id: block.BlockId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get block tasks", slog.String("error", err.Error()))
			return nil, err
		}

		taskModels := make([]eventModel.GetUserTaskAnswer, 0)

		for _, task := range tasks.Tasks {
			taskInfo, err := s.eventClient.GetEditorUserStatsTask(ctx, &dto.UserEventIds{UserId: req.UserId, EventId: task.TaskId})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get editor user stats task", slog.String("error", err.Error()))
				return nil, err
			}

			if taskInfo == nil {
				s.log.With(slog.String("op", op)).Debug("editor user stats task is nil", slog.String("task", task.TaskId), slog.String("user", req.UserId))
				continue
			}

			taskModels = append(taskModels, eventModel.GetUserTaskAnswer{
				TaskId:      taskInfo.TaskId,
				Name:        taskInfo.Name,
				Type:        taskInfo.Type,
				Status:      taskInfo.Status,
				Options:     eventModel.TaskOptionsFromDto(taskInfo.Options),
				UserAnswers: taskInfo.UserAnswerIds,
				UserPoints:  taskInfo.UserPoints,
				Points:      taskInfo.Points,
			})
		}
		blockModels = append(blockModels, eventModel.TaskBlock{
			BlockId: block.BlockId,
			Name:    block.Name,
			Tasks:   taskModels,
		})
	}

	editorStats.Blocks = blockModels

	s.log.With(slog.String("op", op)).Info("got editor stats", slog.Any("user", req.UserId), slog.Any("event", req.EventId))

	return &editorStats, nil
}
