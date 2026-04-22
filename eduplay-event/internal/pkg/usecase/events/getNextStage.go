package event

import (
	"context"
	"errors"
	"slices"

	// "errors"
	"fmt"
	"log/slog"

	errs "eduplay-event/internal/storage"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) GetNextStage(ctx context.Context, in *dto.UserEventIds) (*dto.NextStageInfo, error) {
	const op = "Events.UseCase.GetNextStage"

	log := a.log.With(
		slog.String("op", op),
	)

	nextStageInfo := &dto.NextStageInfo{}
	nextTaskId := ""
	// nextBlockId := ""

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
			_, err := a.storage.EndMe(ctx, in.UserId, in.EventId)
			if err != nil {
				log.Error("failed to end me", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
				return nil, err
			}
			return nextStageInfo, nil
		}

		nextStageInfo, err = a.GetNextBlock(ctx, log, 0, in)
		if err != nil {
			log.Error("failed to get next block", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil

	}

	if currTaskId == "" {
		nextStageInfo, err = a.GetNextBlockById(ctx, log, currBlockId, in)
		if err != nil {
			log.Error("failed to get next block", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}
		return nextStageInfo, nil
	}

	taskAnswer, err := a.storage.GetTaskAnswer(ctx, currTaskId, in.UserId)
	if err != nil {
		log.Error("failed to get user task answer", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	if startTime.AsTime().Format("02.01.2006 15:04:05.000") != "01.01.1970 00:00:00.000" || taskAnswer == nil {
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

		// _, err = a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
		// 	UserId:  in.UserId,
		// 	EventId: in.EventId,
		// 	BlockId: currBlockId,
		// 	TaskId:  currTaskId,
		// })
		// if err != nil {
		// 	log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		// 	return nil, err
		// }

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

	if currBlock.IsParallel {
		nextStageInfo.Type = "block"
		nextStageBlock := &dto.NextStageBlock{
			BlockId: currBlockId,
			Name:    currBlock.Name,
		}

		currBlockTasks, err := a.storage.GetUserBlockTasksShort(ctx, currBlockId, in.UserId)
		if err != nil {
			log.Error("failed to get block tasks short", err.Error(), slog.String("block", currBlockId))
			return nil, err
		}

		for _, task := range currBlockTasks {
			if !task.IsCompleted {
				nextStageBlock.Tasks = currBlockTasks

				nextStageInfo.Block = nextStageBlock

				_, err = a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
					UserId:  in.UserId,
					EventId: in.EventId,
					BlockId: currBlockId,
					TaskId:  "",
				})
				if err != nil {
					log.Error("failed to put next stage", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
					return nil, err
				}

				return nextStageInfo, nil
			}
		}

		nextStageInfo, err = a.GetNextBlock(ctx, log, currBlock.Order, in)
		if err != nil {
			log.Error("failed to get next block", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	currTask, err := a.storage.GetTaskById(ctx, currTaskId)
	if err != nil {
		log.Error("failed to get task", err.Error(), slog.String("task", currTaskId))
		return nil, err
	}

	currBlockTasks, err := a.storage.GetBlockTasks(ctx, currBlockId)
	if err != nil {
		log.Error("failed to get block tasks", err.Error(), slog.String("block", currBlockId))
		return nil, err
	}

	fmt.Println("currTask.Order ", currTask.Order, " len(currBlockTasks.Tasks) ", len(currBlockTasks.Tasks))
	if currTask.Order != int64(len(currBlockTasks.Tasks)) {
		// TODO GetNextTask
		nextStageInfo.Type = "task"

		nextTaskId = currBlockTasks.Tasks[currTask.Order].TaskId
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
			Options:     nextTask.Options,
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

	return a.GetNextBlock(ctx, log, currBlock.Order, in)
}

// TODO ___________________________________________________________

func (a *UseCase) GetNextBlock(ctx context.Context, log *slog.Logger, currBlockOrder int64, in *dto.UserEventIds) (nextStageInfo *dto.NextStageInfo, err error) {

	fmt.Println("===== func GetNextBlock")
	nextStageInfo = &dto.NextStageInfo{}

	nextBlockOrder := int64(0)

	currEvent, err := a.storage.GetEventBlocks(ctx, in.EventId)
	if err != nil {
		return nil, err
	}

	fmt.Println("currBlockOrder ", currBlockOrder, " len(currEvent.Blocks) ", len(currEvent.Blocks))
	fmt.Println(currBlockOrder == int64(len(currEvent.Blocks)))
	if currBlockOrder == int64(len(currEvent.Blocks)) {
		nextStageInfo.Type = "end"

		fmt.Println(nextStageInfo)
		_, err = a.storage.EndMe(ctx, in.UserId, in.EventId)
		if err != nil {
			return nil, err
		}

		return nextStageInfo, nil
	} else if currBlockOrder == 0 {
		nextBlockOrder = 1

		nextBlockId := currEvent.Blocks[nextBlockOrder-1].BlockId

		return a.GetNextBlockById(ctx, log, nextBlockId, in)
	}

	blockWithConditions, err := a.storage.GetBlockConditionsFull(ctx, currEvent.Blocks[currBlockOrder-1].BlockId)
	if err != nil {
		return nil, err
	}

	userPoints, err := a.storage.GetUserBlockPointsSum(ctx, in.UserId, currEvent.Blocks[currBlockOrder-1].BlockId)
	if err != nil {
		return nil, err
	}

	fmt.Println("userPoints ", userPoints, " len(blockConditions) ", len(blockWithConditions.Conditions))
	blockConditions := blockWithConditions.Conditions

	if len(blockConditions) > 0 {

		isGroupEvent := false
		userGroupInfo := &dto.GetUserGroupOut{}

		event, err := a.storage.GetEvent(ctx, in.EventId)
		if err != nil {
			return nil, err
		}

		if event.GroupEvent {
			isGroupEvent = true
			userGroupInfo, err = a.storage.GetUserGroup(ctx, in.UserId, in.EventId)
			if err != nil {
				return nil, err
			}
		}

		for i := len(blockConditions) - 1; i >= 0; i-- {
			fmt.Println("blockCondition ", i, " is ", blockConditions[i])

			min, max := blockConditions[i].Min, blockConditions[i].Max

			if isGroupEvent && len(blockConditions[i].GroupIds) != 0 && !slices.Contains(blockConditions[i].GroupIds, userGroupInfo.GroupId) {
				continue
			}

			fmt.Println("min ", min, " max ", max)
			if max <= 0 || max < min {
				currBlockTasks, err := a.storage.GetBlockTasks(ctx, currEvent.Blocks[currBlockOrder-1].BlockId)
				if err != nil {
					return nil, err
				}
				for _, task := range currBlockTasks.Tasks {
					max += task.Points
				}
				fmt.Println("counting absolute max bc no max was set in condition ", max)
			}
			fmt.Println("userPoints ", userPoints, " min ", min, " max ", max)
			if userPoints > max {
				// return nil, errors.New("too many points")
				fmt.Println("somehow user points larger than max, so we go to next block")
				break
			}
			if userPoints >= min {
				nextBlockOrder = blockConditions[i].NextBlockOrder

				nextBlockId := currEvent.Blocks[nextBlockOrder-1].BlockId

				nextStageInfo, err = a.GetNextBlockById(ctx, log, nextBlockId, in)
				if err != nil {
					return nil, err
				}

				if nextBlockOrder == currBlockOrder {
					fmt.Println("nextBlockOrder == currBlockOrder", nextBlockOrder, currBlockOrder)
					err = a.storage.ClearBlockAnswers(ctx, in.UserId, nextBlockId)
					if err != nil {
						if errors.Is(err, errs.ErrNotFound) {
							a.log.Debug("Block not completed somehow")
						}
						return nil, err
					}
				}

				return nextStageInfo, nil
			}
		}
		// return nil, errors.New("too few points")

		nextBlockOrder = currBlockOrder + 1
		nextBlockId := currEvent.Blocks[nextBlockOrder-1].BlockId

		nextStageInfo, err = a.GetNextBlockById(ctx, log, nextBlockId, in)
		if err != nil {
			return nil, err
		}

		err = a.storage.ClearBlockAnswers(ctx, in.UserId, nextBlockId)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				a.log.Debug("Block not completed somehow")
			}
			return nil, err
		}

		return nextStageInfo, nil
	}

	nextBlockOrder = currBlockOrder + 1

	nextBlockId := currEvent.Blocks[nextBlockOrder-1].BlockId

	nextStageInfo, err = a.GetNextBlockById(ctx, log, nextBlockId, in)
	if err != nil {
		return nil, err
	}

	return nextStageInfo, nil
}

func (a *UseCase) GetNextBlockById(ctx context.Context, log *slog.Logger, nextBlockId string, in *dto.UserEventIds) (nextStageInfo *dto.NextStageInfo, err error) {
	fmt.Println("===== func GetNextBlockById")
	nextStageInfo = &dto.NextStageInfo{}

	nextTaskId := ""

	fmt.Println("---- ", nextBlockId)

	nextBlock, err := a.storage.GetBlockInfo(ctx, nextBlockId)
	if err != nil {
		return nil, err
	}
	nextBlockTasks, err := a.storage.GetUserBlockTasksShort(ctx, nextBlockId, in.UserId)
	if err != nil {
		return nil, err
	}

	if nextBlock.IsParallel {
		fmt.Println("=========== ", nextBlock.IsParallel)
		nextStageInfo.Type = "block"
		nextStageBlock := &dto.NextStageBlock{
			BlockId:    nextBlockId,
			Name:       nextBlock.Name,
			Tasks:      nextBlockTasks,
			IsParallel: true,
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
	if len(nextBlockTasks) == 0 {
		return a.GetNextBlock(ctx, log, nextBlock.Order, in)
	}
	nextTaskId = nextBlockTasks[0].TaskId
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

// func (a *UseCase) GetNextTask(ctx context.Context, сurrTaskId string, currBlockId string, in *dto.UserEventIds) (nextStageInfo *dto.NextStageInfo, err error) {
// 	fmt.Println("===== func GetNextTask")

// 	nextStageInfo = &dto.NextStageInfo{}

// 	currTask, err := a.storage.GetTaskById(ctx, сurrTaskId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	currBlock, err := a.storage.GetBlockInfo(ctx, currBlockId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	currBlockTasks, err := a.storage.GetBlockTasks(ctx, currBlockId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if currBlock.IsParallel {
// 		nextStageInfo.Type = "block"
// 		nextStageBlock := &dto.NextStageBlock{
// 			BlockId: currBlockId,
// 			Name:    currBlock.Name,
// 		}
// 		for _, task := range currBlockTasks.Tasks {
// 			nextStageBlock.Tasks = append(nextStageBlock.Tasks, &dto.NextStageTaskShort{
// 				TaskId:      task.TaskId,
// 				Name:        task.Name,
// 				Time:        task.Time,
// 				IsCompleted: false,
// 			})
// 		}
// 		nextStageInfo.Block = nextStageBlock

// 		_, err := a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
// 			UserId:  in.UserId,
// 			EventId: in.EventId,
// 			BlockId: currBlockId,
// 			TaskId:  "",
// 		})
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nextStageInfo, nil
// 	}

// 	if currTask.Order == int64(len(currBlockTasks.Tasks)) {
// 		return a.GetNextBlock(ctx, currBlock.Order, in)
// 	}

// 	nextTaskId := currBlockTasks.Tasks[currTask.Order+1].TaskId
// 	nextTask, err := a.storage.GetTaskById(ctx, nextTaskId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nextStageInfo.Task = &dto.NextStageTask{
// 		TaskId:      nextTaskId,
// 		BlockId:     currBlockId,
// 		Name:        nextTask.Name,
// 		Description: nextTask.Description,
// 		Type:        nextTask.Type,
// 		Options:     nextTask.Options,
// 		Files:       nextTask.Files,
// 		Time:        nextTask.Time,
// 		Timestamp:   nil,
// 	}

// 	_, err = a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
// 		UserId:  in.UserId,
// 		EventId: in.EventId,
// 		BlockId: currBlockId,
// 		TaskId:  nextTaskId,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nextStageInfo, nil
// }
