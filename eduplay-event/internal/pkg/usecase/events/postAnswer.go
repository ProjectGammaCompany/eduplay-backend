package event

import (
	"context"
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

	corrAnswerIds, corrAnswers, allOptions := GetCorrectAnswers(task)

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
			AnswerIds:   in.Answer,
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
			Answer:      make([]string, 0),
			AnswerIds:   in.Answer,
			Points:      0,
			Status:      "",
			RightAnswer: corrAnswers,
		}
		for i, answer := range corrAnswerIds {
			if in.Answer[0] == answer {
				ans.Points = task.Points
				ans.Status = "correct"
				ans.Answer = append(ans.Answer, corrAnswers[i])

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
			Answer:      make([]string, 0),
			AnswerIds:   in.Answer,
			Points:      0,
			Status:      "",
			RightAnswer: corrAnswers,
		}

		count := 0

		for _, userAnswer := range in.Answer {
			ans.Answer = append(ans.Answer, allOptions[userAnswer])
			for _, correctAnswerId := range corrAnswerIds {
				if userAnswer == correctAnswerId {
					count++
				}
			}
		}

		ans.Points = int64(count) * task.Points / int64(len(corrAnswers))

		log.Info("checking answer", slog.Int("count", count), slog.Int("len", len(corrAnswers)), slog.Int64("points", ans.Points))

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
	case 3:
		ans := &dto.Answer{
			TaskId:      in.TaskId,
			UserId:      in.UserId,
			Answer:      in.Answer,
			AnswerIds:   make([]string, 0),
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

		return ans, nil

	default:
		ans := &dto.Answer{
			TaskId:      in.TaskId,
			UserId:      in.UserId,
			Answer:      make([]string, 0),
			AnswerIds:   in.Answer,
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

		return ans, nil
	}
}

func GetCorrectAnswers(task *dto.Task) ([]string, []string, map[string]string) {
	corrIds := make([]string, 0)
	corr := make([]string, 0)
	options := map[string]string{}

	for _, option := range task.Options {
		options[option.OptionId] = option.Value
		if option.IsCorrect {
			corrIds = append(corrIds, option.OptionId)
			corr = append(corr, option.Value)
		}
	}

	return corrIds, corr, options
}
