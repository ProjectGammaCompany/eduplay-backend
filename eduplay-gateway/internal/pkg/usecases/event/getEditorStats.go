package event

import (
	"context"
	dto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetEditorStats(ctx context.Context, req *eventModel.Id) (*eventModel.EditorStats, error) {
	const op = "event.UseCase.GetEditorStats"

	s.log.With(slog.String("op", op)).Info("attempting to get editor stats")

	editorStats := eventModel.EditorStats{}
	groupsDtos := make([]eventModel.GroupDTO, 0)

	event, err := s.eventClient.GetEvent(ctx, &dto.Id{Id: req.Id})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return nil, err
	}

	editorStats.GroupEvent = event.GroupEvent

	if event.GroupEvent {
		groups, err := s.eventClient.GetGroups(ctx, &dto.Id{Id: req.Id})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get groups", slog.String("error", err.Error()))
			return nil, err
		}

		for _, group := range groups.Groups {
			userDtos := make([]eventModel.UserDTO, 0)

			users, err := s.eventClient.GetGroupUsers(ctx, &dto.Id{Id: group.Id})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get group users", slog.String("error", err.Error()))
				return nil, err
			}

			for _, user := range users.Users {
				userInfo, err := s.userClient.GetProfile(ctx, user.Id)
				if err != nil {
					s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
					return nil, err
				}
				userAnswers, err := s.eventClient.GetUserAnswers(ctx, &dto.UserEventIds{EventId: req.Id, UserId: user.Id})
				if err != nil {
					s.log.With(slog.String("op", op)).Error("failed to get user answers", slog.String("error", err.Error()))
					return nil, err
				}

				userDtos = append(userDtos, eventModel.UserDTO{
					UserId:   user.Id,
					Username: userInfo.UserName,
					Avatar:   userInfo.Avatar,
					Points:   userAnswers.Points,
					Answers: eventModel.UserAnswers{
						Correct: userAnswers.Correct,
						Total:   userAnswers.Total,
					},
				})
			}

			groupsDtos = append(groupsDtos, eventModel.GroupDTO{
				GroupId: group.Id,
				Name:    group.Login,
				Users:   userDtos,
			})
		}

		editorStats.Groups = groupsDtos
	} else {
		userDtos := make([]eventModel.UserDTO, 0)

		userStats, err := s.eventClient.GetEventUsers(ctx, &dto.Id{Id: req.Id})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user stats", slog.String("error", err.Error()))
			return nil, err
		}

		for _, user := range userStats.Users {
			userInfo, err := s.userClient.GetProfile(ctx, user.Id)
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
				return nil, err
			}
			userAnswers, err := s.eventClient.GetUserAnswers(ctx, &dto.UserEventIds{EventId: req.Id, UserId: user.Id})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get user answers", slog.String("error", err.Error()))
				return nil, err
			}

			userDtos = append(userDtos, eventModel.UserDTO{
				UserId:   user.Id,
				Username: userInfo.UserName,
				Avatar:   userInfo.Avatar,
				Points:   userAnswers.Points,
				Answers: eventModel.UserAnswers{
					Correct: userAnswers.Correct,
					Total:   userAnswers.Total,
				},
			})
		}

		editorStats.Users = userDtos
	}

	s.log.With(slog.String("op", op)).Info("got editor stats", slog.Any("event", req.Id))

	return &editorStats, nil
}
