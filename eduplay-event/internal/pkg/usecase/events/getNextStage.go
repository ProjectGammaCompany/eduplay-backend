package event

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetNextStage(ctx context.Context, in *dto.UserEventIds) (*dto.NextStageInfo, error) {
	const op = "Events.UseCase.GetNextStage"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting next stage")

	nextStageInfo := &dto.NextStageInfo{}
	nextTaskId := ""
	nextBlockId := ""

	log.Info("getting next stage")

	_, currTaskId, currBlockId, finished, startTime, err := a.storage.GetNextStage(ctx, in)
	if err != nil {
		log.Error("failed to get next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	fmt.Println("gotNextStage: ", currTaskId, " ", currBlockId, " ", finished, " ", startTime, " ", in.UserId, " ", in.EventId)

	if finished {
		nextStageInfo.Type = "end"
		return nextStageInfo, nil
	}

	if currBlockId == "" {

		fmt.Println("currBlockId is empty")

		eventBlocks, err := a.storage.GetEventBlocks(ctx, in.EventId)
		if err != nil {
			log.Error("failed to get event blocks", err.Error(), slog.String("event", in.EventId))
			return nil, err
		}

		if len(eventBlocks.Blocks) == 0 {
			nextStageInfo.Type = "end"
			return nextStageInfo, nil
		}

		nextStageInfo, err = a.GetNextBlock(ctx, 0, in)
		if err != nil {
			log.Error("failed to get next block", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil

	}

	if currTaskId == "" {
		nextStageInfo, err = a.GetNextBlockById(ctx, currBlockId, in)
		if err != nil {
			log.Error("failed to get next block", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}
		return nextStageInfo, nil
	}

	if startTime.AsTime().Format("02.01.2006 15:04:05.000") != "01.01.1970 00:00:00.000" {
		fmt.Println("startTime ", startTime.AsTime().Format("02.01.2006 15:04:05.000"))
		nextStageInfo.Type = "task"
		currTask, err := a.storage.GetTaskById(ctx, currTaskId)
		if err != nil {
			log.Error("failed to get task", err.Error(), slog.String("task", currTaskId))
			return nil, err
		}

		nextStageTask := &dto.NextStageTask{
			TaskId:      currTaskId,
			BlockId:     currBlockId,
			Name:        currTask.Name,
			Description: currTask.Description,
			Type:        currTask.Type,
			Options:     currTask.Options,
			Files:       currTask.Files,
			Time:        currTask.Time,
			Timestamp:   startTime,
		}
		nextStageInfo.Task = nextStageTask

		_, err = a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
			UserId:  in.UserId,
			EventId: in.EventId,
			BlockId: currBlockId,
			TaskId:  currTaskId,
		})
		if err != nil {
			log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	// currTask, err := a.storage.GetTaskById(ctx, currTaskId)
	// if err != nil {
	// 	log.Error("failed to get task", err.Error(), slog.String("task", currTaskId))
	// 	return nil, err
	// }

	currBlock, err := a.storage.GetBlockInfo(ctx, currBlockId)
	if err != nil {
		log.Error("failed to get block", err.Error(), slog.String("block", currBlockId))
		return nil, err
	}

	currBlockTasks, err := a.storage.GetBlockTasks(ctx, currBlockId)
	if err != nil {
		log.Error("failed to get block tasks", err.Error(), slog.String("block", currBlockId))
		return nil, err
	}

	if currBlock.IsParallel {
		nextStageInfo.Type = "block"
		nextStageBlock := &dto.NextStageBlock{
			BlockId: nextBlockId,
			Name:    currBlock.Name,
		}
		for _, task := range currBlockTasks.Tasks {
			nextStageBlock.Tasks = append(nextStageBlock.Tasks, &dto.NextStageTaskShort{
				TaskId:      task.TaskId,
				Name:        task.Name,
				Time:        task.Time,
				IsCompleted: false,
			})
		}
		nextStageInfo.Block = nextStageBlock

		_, err := a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
			UserId:  in.UserId,
			EventId: in.EventId,
			BlockId: nextBlockId,
			TaskId:  "",
		})
		if err != nil {
			log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	currTask, err := a.storage.GetTaskById(ctx, currTaskId)
	if err != nil {
		log.Error("failed to get task", err.Error(), slog.String("task", currTaskId))
		return nil, err
	}

	fmt.Println("currTask.Order ", currTask.Order, " len(currBlockTasks.Tasks) ", len(currBlockTasks.Tasks))
	if currTask.Order != int64(len(currBlockTasks.Tasks)) {
		nextStageInfo.Type = "task"

		nextTaskId = currBlockTasks.Tasks[currTask.Order+1].TaskId
		fmt.Println("currTaskId ", currTaskId, " nextTaskId ", nextTaskId)
		nextTask, err := a.storage.GetTaskById(ctx, nextTaskId)
		if err != nil {
			log.Error("failed to get task", err.Error(), slog.String("task", nextTaskId))
			return nil, err
		}

		nextStageInfo.Task = &dto.NextStageTask{
			TaskId:      nextTaskId,
			BlockId:     nextTask.BlockId,
			Name:        nextTask.Name,
			Description: nextTask.Description,
			Type:        nextTask.Type,
			Files:       nextTask.Files,
			Time:        nextTask.Time,
			Timestamp:   startTime,
		}

		_, err = a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
			UserId:  in.UserId,
			EventId: in.EventId,
			BlockId: nextTask.BlockId,
			TaskId:  nextTaskId,
		})
		if err != nil {
			log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	eventBlocks, err := a.storage.GetEventBlocks(ctx, in.EventId)
	if err != nil {
		log.Error("failed to get event blocks", err.Error(), slog.String("event", in.EventId))
		return nil, err
	}

	if currBlock.Order == int64(len(eventBlocks.Blocks)) {
		nextStageInfo.Type = "end"

		_, err = a.storage.EndMe(ctx, in.UserId, in.EventId)
		if err != nil {
			log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	nextBlockId = eventBlocks.Blocks[currBlock.Order+1].BlockId

	nextBlock, err := a.storage.GetBlockInfo(ctx, nextBlockId)
	if err != nil {
		log.Error("failed to get block", err.Error(), slog.String("block", nextBlockId))
		return nil, err
	}

	if nextBlock.IsParallel {
		nextStageInfo.Type = "block"
		nextStageBlock := &dto.NextStageBlock{
			BlockId: nextBlockId,
			Name:    nextBlock.Name,
		}
		for _, task := range currBlockTasks.Tasks {
			nextStageBlock.Tasks = append(nextStageBlock.Tasks, &dto.NextStageTaskShort{
				TaskId:      task.TaskId,
				Name:        task.Name,
				Time:        task.Time,
				IsCompleted: false,
			})
		}
		nextStageInfo.Block = nextStageBlock

		_, err := a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
			UserId:  in.UserId,
			EventId: in.EventId,
			BlockId: nextBlockId,
			TaskId:  "",
		})
		if err != nil {
			log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	nextTaskId = currBlockTasks.Tasks[0].TaskId

	nextTask, err := a.storage.GetTaskById(ctx, nextTaskId)
	if err != nil {
		log.Error("failed to get task", err.Error(), slog.String("task", nextTaskId))
		return nil, err
	}

	nextStageInfo.Type = "task"
	nextStageInfo.Task = &dto.NextStageTask{
		TaskId:      nextTaskId,
		BlockId:     nextTask.BlockId,
		Name:        nextTask.Name,
		Description: nextTask.Description,
		Type:        nextTask.Type,
		Files:       nextTask.Files,
		Time:        nextTask.Time,
		Timestamp:   startTime,
	}

	_, err = a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
		UserId:  in.UserId,
		EventId: in.EventId,
		BlockId: nextTask.BlockId,
		TaskId:  nextTaskId,
	})
	if err != nil {
		log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return nextStageInfo, nil
	// return &dto.NextStageInfo{Type: "end"}, nil
}

func (a *UseCase) GetNextBlock(ctx context.Context, currBlockOrder int64, in *dto.UserEventIds) (nextStageInfo *dto.NextStageInfo, err error) {
	fmt.Println("===== func GetNextBlock")

	nextBlockOrder := int64(0)

	currEvent, err := a.storage.GetEventBlocks(ctx, in.EventId)
	if err != nil {
		return nil, err
	}

	if currBlockOrder == 0 {
		nextBlockOrder = 1

		nextBlockId := currEvent.Blocks[nextBlockOrder-1].BlockId

		return a.GetNextBlockById(ctx, nextBlockId, in)
	} else if currBlockOrder == int64(len(currEvent.Blocks)) {
		nextStageInfo.Type = "end"

		_, err = a.storage.EndMe(ctx, in.UserId, in.EventId)
		if err != nil {
			return nil, err
		}

		return nextStageInfo, nil
	}
	blockWithConditions, err := a.storage.GetBlockConditionsFull(ctx, currEvent.Blocks[currBlockOrder-1].BlockId)
	if err != nil {
		return nil, err
	}

	userPoints, err := a.storage.GetUserBlockPointsSum(ctx, in.UserId, currEvent.Blocks[currBlockOrder-1].BlockId)
	if err != nil {
		return nil, err
	}

	blockConditions := blockWithConditions.Conditions

	if len(blockConditions) > 0 {
		for i := len(blockConditions) - 1; i > 0; i-- {
			min, max := blockConditions[i].Min, blockConditions[i].Max
			if max == 0 {
				currBlockTasks, err := a.storage.GetBlockTasks(ctx, currEvent.Blocks[currBlockOrder-1].BlockId)
				if err != nil {
					return nil, err
				}
				for _, task := range currBlockTasks.Tasks {
					max += task.Points
				}
			}
			if userPoints > max {
				return nil, errors.New("too many points")
			}
			if userPoints > min {
				nextBlockOrder = blockConditions[i].NextBlockOrder

				nextBlockId := currEvent.Blocks[nextBlockOrder-1].BlockId

				nextStageInfo, err = a.GetNextBlockById(ctx, nextBlockId, in)
				if err != nil {
					return nil, err
				}

				return nextStageInfo, nil
			}
		}
		return nil, errors.New("too few points")
	}

	nextBlockOrder = currBlockOrder + 1

	nextBlockId := currEvent.Blocks[nextBlockOrder-1].BlockId

	nextStageInfo, err = a.GetNextBlockById(ctx, nextBlockId, in)
	if err != nil {
		return nil, err
	}

	return nextStageInfo, nil
}

func (a *UseCase) GetNextBlockById(ctx context.Context, nextBlockId string, in *dto.UserEventIds) (nextStageInfo *dto.NextStageInfo, err error) {
	fmt.Println("===== func GetNextBlockById")
	nextStageInfo = &dto.NextStageInfo{}

	nextTaskId := ""

	fmt.Println("---- ", nextBlockId)

	nextBlock, err := a.storage.GetBlockInfo(ctx, nextBlockId)
	if err != nil {
		return nil, err
	}
	nextBlockTasks, err := a.storage.GetBlockTasks(ctx, nextBlockId)
	if err != nil {
		return nil, err
	}

	if nextBlock.IsParallel {
		nextStageInfo.Type = "block"
		nextStageBlock := &dto.NextStageBlock{
			BlockId: nextBlockId,
			Name:    nextBlock.Name,
		}
		fmt.Println("=========== ", nextBlock.IsParallel)
		for _, task := range nextBlockTasks.Tasks {
			nextStageBlock.Tasks = append(nextStageBlock.Tasks, &dto.NextStageTaskShort{
				TaskId:      task.TaskId,
				Name:        task.Name,
				Time:        task.Time,
				IsCompleted: false,
			})
		}
		nextStageInfo.Block = nextStageBlock

		_, err := a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
			UserId:  in.UserId,
			EventId: in.EventId,
			BlockId: nextBlockId,
			TaskId:  "",
		})
		if err != nil {
			return nil, err
		}

		return nextStageInfo, nil
	}

	fmt.Println("===== get block task")

	nextStageInfo.Type = "task"
	if len(nextBlockTasks.Tasks) == 0 {
		return a.GetNextBlock(ctx, nextBlock.Order, in)
	}
	nextTaskId = nextBlockTasks.Tasks[0].TaskId
	nextTask, err := a.storage.GetTaskById(ctx, nextTaskId)
	if err != nil {
		return nil, err
	}

	nextStageInfo.Task = &dto.NextStageTask{
		TaskId:      nextTaskId,
		BlockId:     nextBlockId,
		Name:        nextTask.Name,
		Description: nextTask.Description,
		Type:        nextTask.Type,
		Options:     nextTask.Options,
		Files:       nextTask.Files,
		Time:        nextTask.Time,
		Timestamp:   nil,
	}

	_, err = a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
		UserId:  in.UserId,
		EventId: in.EventId,
		BlockId: nextBlockId,
		TaskId:  nextTaskId,
	})
	if err != nil {
		return nil, err
	}

	return nextStageInfo, nil
}

// func (a *UseCase) GetNextTask(ctx context.Context, currTaskOrder string, in *dto.UserEventIds) (nextStageInfo *dto.NextStageInfo, err error) {
// 	fmt.Println("===== func GetNextTask")

// }
