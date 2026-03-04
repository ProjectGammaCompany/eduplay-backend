package event

import (
	"context"
	"fmt"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostAnswer(ctx context.Context, in *dto.Answer) (*dto.Answer, error) {
	const op = "Events.UseCase.PostAnswer"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("getting task by task id")

	task, err := a.storage.GetTaskById(ctx, in.TaskId)
	if err != nil {
		log.Error("failed to get task by task id", err.Error(), slog.String("block", in.TaskId))
		return nil, err
	}

	corrAnswers := GetCorrectAnswers(task)

	// TODO: put empty timestamp
	_, err = a.storage.PutTimestamp(ctx, in.UserId, in.EventId, nil)
	if err != nil {
		return nil, err
	}

	log.Info("checking answer")

	switch task.Type {
	case 0:
		ans := &dto.Answer{
			TaskId:      in.TaskId,
			UserId:      in.UserId,
			Answer:      in.Answer,
			Points:      task.Points,
			Status:      "correct",
			RightAnswer: corrAnswers,
		}

		_, err := a.storage.PostAnswer(ctx, ans)
		if err != nil {
			return nil, err
		}

		return ans, nil
	case 1:
		ans := &dto.Answer{
			TaskId:      in.TaskId,
			UserId:      in.UserId,
			Answer:      in.Answer,
			Points:      0,
			Status:      "",
			RightAnswer: corrAnswers,
		}
		for _, answer := range corrAnswers {
			if in.Answer[0] == answer {
				ans.Points = task.Points
				ans.Status = "correct"

				_, err := a.storage.PostAnswer(ctx, ans)
				if err != nil {
					return nil, err
				}

				return ans, nil
			}
		}
		ans.Status = "incorrect"

		_, err := a.storage.PostAnswer(ctx, ans)
		if err != nil {
			return nil, err
		}

		return ans, nil
	case 2:
		ans := &dto.Answer{
			TaskId:      in.TaskId,
			UserId:      in.UserId,
			Answer:      in.Answer,
			Points:      0,
			Status:      "",
			RightAnswer: corrAnswers,
		}

		count := 0

		for _, userAnswer := range in.Answer {
			for _, correctAnswer := range corrAnswers {
				if userAnswer == correctAnswer {
					count++
				}
			}
		}

		ans.Points = int64(count) * task.Points / int64(len(task.Options))
		fmt.Println(count)

		if count == len(corrAnswers) {
			ans.Status = "correct"
		} else if count > 0 {
			ans.Status = "partial"
		} else {
			ans.Status = "incorrect"
		}

		_, err := a.storage.PostAnswer(ctx, ans)
		if err != nil {
			return nil, err
		}

		return ans, nil
	default:
		ans := &dto.Answer{
			TaskId:      in.TaskId,
			UserId:      in.UserId,
			Answer:      in.Answer,
			Points:      0,
			Status:      "",
			RightAnswer: corrAnswers,
		}
		if in.Answer[0] == corrAnswers[0] {
			ans.Points = task.Points
			ans.Status = "correct"
		} else {
			ans.Points = 0
			ans.Status = "incorrect"
		}

		_, err := a.storage.PostAnswer(ctx, ans)
		if err != nil {
			return nil, err
		}

		block, err := a.storage.GetBlockInfo(ctx, task.BlockId)
		if err != nil {
			return nil, err
		}

		event, err := a.storage.GetEvent(ctx, block.EventId)
		if err != nil {
			return nil, err
		}

		_, err = a.storage.PutTimestamp(ctx, in.UserId, event.EventId, nil)
		if err != nil {
			return nil, err
		}

		return ans, nil
	}
}

func GetCorrectAnswers(task *dto.Task) []string {
	corr := make([]string, 0)

	for i, option := range task.Options {
		if option.IsCorrect {
			corr = append(corr, task.Options[i].Value)
		}
	}

	return corr
}
