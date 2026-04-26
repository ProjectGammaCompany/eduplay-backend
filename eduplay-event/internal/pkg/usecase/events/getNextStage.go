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
		if errors.Is(err, errs.ErrNotFound) {
			nextStageInfo.Type = errs.ErrNotFound.Error()
			return nextStageInfo, nil
		}
		log.Error("failed to get next stage", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	log.Debug("got next stage", slog.String("currTaskId", currTaskId), slog.String("currBlockId", currBlockId), slog.Bool("finished", finished), slog.String("startTime", startTime.String()), slog.String("userId", in.UserId), slog.String("eventId", in.EventId))

	if finished {
		nextStageInfo.Type = "end"
		return nextStageInfo, nil
	}

	if currBlockId == "" {
		nextStageInfo, err = a.HandleNoBlock(ctx, log, nextStageInfo, in)
		if err != nil {
			log.Error("failed to handle no block", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}
		return nextStageInfo, nil
	}

	if currTaskId == "" {
		nextStageInfo, err = a.GetNextBlockById(ctx, log, currBlockId, in)
		if err != nil {
			log.Error("failed to get next block", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}
		return nextStageInfo, nil
	}

	taskAnswer, err := a.storage.GetTaskAnswer(ctx, currTaskId, in.UserId)
	if err != nil {
		log.Error("failed to get user task answer", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	if startTime.AsTime().IsZero() || taskAnswer == nil {
		log.Debug("return same task ", slog.String("startTime", startTime.AsTime().Format("02.01.2006 15:04:05.000")))
		nextStageInfo.Type = "task"
		currTask, err := a.storage.GetTaskById(ctx, currTaskId)
		if err != nil {
			log.Error("failed to get task", slog.Any("error", err), slog.String("task", currTaskId))
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

		return nextStageInfo, nil
	}

	currBlock, err := a.storage.GetBlockInfo(ctx, currBlockId)
	if err != nil {
		log.Error("failed to get block", slog.Any("error", err), slog.String("block", currBlockId))
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
			log.Error("failed to get block tasks short", slog.Any("error", err), slog.String("block", currBlockId))
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
					log.Error("failed to put next stage", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
					return nil, err
				}

				return nextStageInfo, nil
			}
		}

		nextStageInfo, err = a.GetNextBlock(ctx, log, currBlock.Order, in)
		if err != nil {
			log.Error("failed to get next block", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	currTask, err := a.storage.GetTaskById(ctx, currTaskId)
	if err != nil {
		log.Error("failed to get task", slog.Any("error", err), slog.String("task", currTaskId))
		return nil, err
	}

	currBlockTasks, err := a.storage.GetUserBlockTasksShort(ctx, currBlockId, in.UserId)
	if err != nil {
		log.Error("failed to get block tasks", slog.Any("error", err), slog.String("block", currBlockId))
		return nil, err
	}

	fmt.Println("currTask.Order ", currTask.Order, " len(currBlockTasks.Tasks) ", len(currBlockTasks))
	if currTask.Order != int64(len(currBlockTasks)) {
		// TODO GetNextTask
		nextStageInfo.Type = "task"

		nextTaskId = currBlockTasks[currTask.Order].TaskId
		fmt.Println("currTaskId ", currTaskId, " nextTaskId ", nextTaskId)
		nextTask, err := a.storage.GetTaskById(ctx, nextTaskId)
		if err != nil {
			log.Error("failed to get task", slog.Any("error", err), slog.String("task", nextTaskId))
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
			log.Error("failed to put next stage", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}

		return nextStageInfo, nil
	}

	return a.GetNextBlock(ctx, log, currBlock.Order, in)
}

// TODO ___________________________________________________________

func (a *UseCase) HandleNoBlock(ctx context.Context, log *slog.Logger, nextStageInfo *dto.NextStageInfo, in *dto.UserEventIds) (*dto.NextStageInfo, error) {
	log.Debug("currBlockId is empty")

	eventBlocks, err := a.storage.GetEventBlocks(ctx, in.EventId)
	if err != nil {
		log.Error("failed to get event blocks", slog.Any("error", err), slog.String("event", in.EventId))
		return nil, err
	}

	if len(eventBlocks.Blocks) == 0 {
		nextStageInfo.Type = "end"
		_, err := a.storage.EndMe(ctx, in.UserId, in.EventId)
		if err != nil {
			log.Error("failed to end me", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}
		return nextStageInfo, nil
	}

	nextStageInfo, err = a.GetNextBlock(ctx, log, 0, in)
	if err != nil {
		log.Error("failed to get next block", slog.Any("error", err), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return nextStageInfo, nil
}

func (a *UseCase) GetNextBlock(ctx context.Context, log *slog.Logger, currBlockOrder int64, in *dto.UserEventIds) (nextStageInfo *dto.NextStageInfo, err error) {

	log.Debug("===== func GetNextBlock")
	nextStageInfo = &dto.NextStageInfo{}

	nextBlockOrder := int64(0)

	currEvent, err := a.storage.GetEventBlocks(ctx, in.EventId)
	if err != nil {
		return nil, err
	}

	log.Debug("block order", slog.Int64("currBlockOrder ", currBlockOrder), slog.Int64("len(currEvent.Blocks) ", int64(len(currEvent.Blocks))))
	log.Debug("is block last in event?", slog.Bool("currBlockOrder == len(currEvent.Blocks)", currBlockOrder == int64(len(currEvent.Blocks))))
	if currBlockOrder == int64(len(currEvent.Blocks)) {
		nextStageInfo.Type = "end"

		log.Debug("block is last", slog.Any("nextStageInfo", nextStageInfo))
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

	log.Debug("check for conditions", slog.Int64("userPoints", userPoints), slog.Int("len(blockConditions) ", len(blockWithConditions.Conditions)))
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
			log.Debug("curr condition", slog.Int("block condition", i), slog.Int("len(blockConditions)", len(blockConditions)))

			min, max := blockConditions[i].Min, blockConditions[i].Max

			if isGroupEvent && len(blockConditions[i].GroupIds) != 0 && !slices.Contains(blockConditions[i].GroupIds, userGroupInfo.GroupId) {
				continue
			}

			fmt.Println("points condition", userPoints, min, max)
			log.Debug("points condition", slog.Int64("userPoints ", userPoints), slog.Int64("min ", min), slog.Int64("max ", max))
			if max <= 0 || max < min {
				max, err = a.storage.GetBlockMaxPoints(ctx, currEvent.Blocks[currBlockOrder-1].BlockId)
				if err != nil {
					return nil, err
				}
				log.Debug("counting absolute max bc no max was set in condition", slog.Int64("max ", max))
			}
			log.Debug("points condition", slog.Int64("userPoints ", userPoints), slog.Int64("min ", min), slog.Int64("max ", max))
			// if userPoints > max {
			// 	// return nil, errors.New("too many points")
			// 	log.Warn("somehow user points larger than max, so we go to next block")
			// 	break
			// }
			if userPoints >= min && userPoints <= max {
				nextBlockOrder = blockConditions[i].NextBlockOrder

				fmt.Println("nextBlockOrder", nextBlockOrder)
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
