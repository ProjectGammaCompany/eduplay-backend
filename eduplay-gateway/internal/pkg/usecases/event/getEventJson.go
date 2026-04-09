package event

import (
	"context"
	"crypto/sha256"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"encoding/hex"
	"log/slog"
)

func (s *UseCase) GetEventJson(ctx context.Context, eventId string) (*eventModel.EventDownloadFull, error) {
	const op = "event.UseCase.GetEventJson"

	s.log.With(slog.String("op", op)).Info("attempting to get event json")

	eventJson := &eventModel.EventDownloadFull{}
	eventJson.Files = make([]string, 0)
	eventGroups := make([]eventModel.GroupDownload, 0)

	uniqueFiles := make(map[string]bool)

	event, err := s.eventClient.GetEvent(ctx, &eventDto.Id{Id: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return nil, err
	}

	eventJson.EventDownload.AuthorId = append(eventJson.EventDownload.AuthorId, event.OwnerId)
	eventJson.EventDownload.EventId = event.EventId
	eventJson.EventDownload.Title = event.Title
	eventJson.EventDownload.Description = event.Description
	eventJson.EventDownload.Tags = event.Tags
	eventJson.EventDownload.Cover = event.Cover
	eventJson.EventDownload.StartDate = event.StartDate.AsTime().Format("02.01.2006 15:04:05.000")
	eventJson.EventDownload.EndDate = event.EndDate.AsTime().Format("02.01.2006 15:04:05.000")
	eventJson.EventDownload.LastEditionDate = event.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000")
	eventJson.EventDownload.GroupEvent = event.GroupEvent

	uniqueFiles[event.Cover] = true

	if event.GroupEvent {
		groups, err := s.eventClient.GetGroups(ctx, &eventDto.Id{Id: eventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get groups", slog.String("error", err.Error()))
			return nil, err
		}

		for _, group := range groups.Groups {
			hash := sha256.Sum256([]byte(group.Password))
			eventGroups = append(eventGroups, eventModel.GroupDownload{
				GroupId:  group.Id,
				EventId:  eventId,
				Login:    group.Login,
				Password: hex.EncodeToString(hash[:]),
			})
		}
	}

	collaborators, err := s.eventClient.GetCollaborators(ctx, &eventDto.Id{Id: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get collaborators", slog.String("error", err.Error()))
		return nil, err
	}

	for _, collaborator := range collaborators.Users {
		eventJson.EventDownload.AuthorId = append(eventJson.EventDownload.AuthorId, collaborator.Id)
	}

	s.log.With(slog.String("op", op)).Info("got event json", slog.Any("event", event.EventId))

	eventBlocks := make([]eventModel.BlockDownload, 0)
	eventBlockConditions := make([]eventModel.ConditionDownload, 0)
	eventTasks := make([]eventModel.TaskDownload, 0)
	eventTaskOptions := make([]eventModel.OptionDownload, 0)
	eventTaskAnswers := make([]eventModel.CorrectAnswerDownload, 0)

	blocks, err := s.eventClient.GetEventBlocks(ctx, &eventDto.Id{Id: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event blocks", slog.String("error", err.Error()))
		return nil, err
	}

	for _, block := range blocks.Blocks {
		blockInfo, err := s.eventClient.GetBlockInfo(ctx, &eventDto.Id{Id: block.BlockId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get block info", slog.String("error", err.Error()))
			return nil, err
		}

		eventBlocks = append(eventBlocks, eventModel.BlockDownload{
			BlockId:       block.BlockId,
			Name:          blockInfo.Name,
			BlockOrder:    blockInfo.Order,
			IsParallel:    blockInfo.IsParallel,
			ShowPoints:    blockInfo.ShowPoints,
			ShowAnswers:   blockInfo.ShowAnswers,
			PartialPoints: blockInfo.PartialPoints,
			EventId:       eventId,
		})

		for _, condition := range block.Conditions {
			eventBlockConditions = append(eventBlockConditions, eventModel.ConditionDownload{
				ConditionId: condition.ConditionId,
				PrevBlockId: condition.PreviousBlockId,
				NextBlockId: condition.NextBlockId,
				GroupName:   condition.GroupIds,
				Min:         &condition.Min,
				Max:         &condition.Max,
			})
		}

		tasks, err := s.eventClient.GetBlockTasks(ctx, &eventDto.Id{Id: block.BlockId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get block tasks", slog.String("error", err.Error()))
			return nil, err
		}

		for _, task := range tasks.Tasks {
			taskInfo, err := s.eventClient.GetTaskById(ctx, &eventDto.Id{Id: task.TaskId})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get task info", slog.String("error", err.Error()))
				return nil, err
			}

			eventTaskFiles := make([]string, len(taskInfo.Files))
			for _, file := range task.Files {
				eventTaskFiles = append(eventTaskFiles, file.Url)
				uniqueFiles[file.Url] = true
			}

			eventTasks = append(eventTasks, eventModel.TaskDownload{
				TaskId:        task.TaskId,
				BlockId:       block.BlockId,
				Name:          taskInfo.Name,
				Description:   taskInfo.Description,
				TaskType:      taskInfo.Type,
				Points:        taskInfo.Points,
				PartialPoints: taskInfo.PartialPoints,
				Time:          taskInfo.Time,
				Order:         taskInfo.Order,
				Files:         eventTaskFiles,
			})

			taskAnswers := make([]string, 0)

			for _, option := range task.Options {
				if option.IsCorrect {
					if task.Type == 3 || task.Type == 4 {
						hash := sha256.Sum256([]byte(option.Value))
						taskAnswers = append(taskAnswers, hex.EncodeToString(hash[:]))
					} else {
						hash := sha256.Sum256([]byte(option.OptionId))
						taskAnswers = append(taskAnswers, hex.EncodeToString(hash[:]))
					}
				}
				eventTaskOptions = append(eventTaskOptions, eventModel.OptionDownload{
					OptionId: option.OptionId,
					TaskId:   task.TaskId,
					Value:    option.Value,
				})
			}

			eventTaskAnswers = append(eventTaskAnswers, eventModel.CorrectAnswerDownload{
				TaskId: task.TaskId,
				Values: taskAnswers,
			})
		}
	}

	for file := range uniqueFiles {
		eventJson.Files = append(eventJson.Files, file)
	}

	eventJson.GroupsDownload = eventGroups
	eventJson.BlocksDownload = eventBlocks
	eventJson.ConditionsDownload = eventBlockConditions
	eventJson.TasksDownload = eventTasks
	eventJson.OptionsDownload = eventTaskOptions
	eventJson.CorrectAnswersDownload = eventTaskAnswers

	return eventJson, nil
}
