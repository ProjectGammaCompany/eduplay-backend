package event

import (
	"context"
	"fmt"
	"log/slog"

	dto "eduplay-event/internal/generated"
)

func (a *UseCase) PostAnswerBatch(ctx context.Context, in *dto.AnswerBatch) (*dto.MessageOut, error) {
	const op = "Events.UseCase.PostAnswerBatch"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("updating user status from answer batch")

	_, _, finished, _, err := a.storage.GetEventProgress(ctx, in.UserId, in.EventId)
	if err != nil {
		log.Error("failed to get event progress", err.Error(), slog.String("event", in.EventId))
		return nil, err
	}

	if finished {
		return &dto.MessageOut{Message: "event is finished"}, nil
	}

	fmt.Println(op, in.CurrentTask, in.CurrentBlock, in.IsDone)

	for _, answer := range in.Answers {
		_, err := a.PostAnswer(ctx, &dto.Answer{
			UserId:  in.UserId,
			EventId: in.EventId,
			TaskId:  answer.TaskId,
			Answer:  answer.Options,
		})

		if err != nil {
			log.Error("failed to update user status from answer batch", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
			return nil, err
		}
	}

	// id, err := a.storage.PostAnswerBatch(ctx, in)
	message, err := a.storage.PutNextStage(ctx, &dto.EventBlockTaskUserIds{
		UserId:   in.UserId,
		EventId:  in.EventId,
		TaskId:   in.CurrentTask,
		BlockId:  in.CurrentBlock,
		Finished: in.IsDone,
	})
	if err != nil {
		log.Error("failed to update user status from answer batch", err.Error(), slog.String("event", in.EventId), slog.String("user", in.UserId))
		return nil, err
	}

	return &dto.MessageOut{Message: message}, nil
}
