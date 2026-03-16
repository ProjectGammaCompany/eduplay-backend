package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetPlayerStats(ctx context.Context, in *eventModel.UserEventIds) (*eventModel.PlayerStats, error) {
	const op = "event.UseCase.GetPlayerStats"

	s.log.With(slog.String("op", op)).Info("attempting to get player stats")

	event, err := s.eventClient.GetEvent(ctx, &eventDto.Id{Id: in.EventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return nil, err
	}

	if !event.Rating {
		userStats, err := s.eventClient.GetUserStats(ctx, eventModel.UserEventIdsToDto(in))
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get player stats", slog.String("error", err.Error()))
			return nil, err
		}

		userInfo, err := s.userClient.GetProfile(ctx, userStats.Id)
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
			return nil, err
		}

		ret := &eventModel.PlayerStats{
			FullStats:  false,
			GroupEvent: false,
			Users:      make([]eventModel.UserStats, 0),
			Groups:     make([]eventModel.GroupStats, 0),
		}

		ret.Users = append(ret.Users, eventModel.UserStats{
			UserId:   userStats.Id,
			Username: userInfo.Email,
			Avatar:   userInfo.Avatar,
			Points:   userStats.Points,
		})

		return ret, nil
	}

	// if event.GroupEvent {
	// 	userGroups, err := s.eventClient.GetGroupUsers(ctx, &eventDto.Id{Id: in.EventId})
	// 	if err != nil {
	// 		s.log.With(slog.String("op", op)).Error("failed to get user groups", slog.String("error", err.Error()))
	// 		return nil, err
	// 	}

	// 	ret := &eventModel.PlayerStats{
	// 		FullStats:  true,
	// 		GroupEvent: true,
	// 		Users:      make([]eventModel.UserStats, 0),
	// 		Groups:     make([]eventModel.GroupStats, 0),
	// 	}

	// 	for i, userGroup := range userGroups.Groups {
	// 		ret.Groups = append(ret.Groups, eventModel.GroupStats{
	// 			GroupId: userGroup.GroupId,
	// 			Name:    userGroup.Name,
	// 			Users:   make([]eventModel.UserStats, 0),
	// 		})
	// 		for _, user := range userGroup.Users {
	// 			ret.Groups[i].Users = append(ret.Groups[i].Users, eventModel.UserStats{
	// 				UserId: user.Id,
	// 			})
	// 			userStats, err := s.eventClient.GetUserStats(ctx, &eventDto.UserEventIds{
	// 				EventId: in.EventId,
	// 				UserId:  in.UserId,
	// 			})
	// 			if err != nil {
	// 				s.log.With(slog.String("op", op)).Error("failed to get user stats", slog.String("error", err.Error()))
	// 				continue
	// 			}
	// 			userInfo, err := s.userClient.GetProfile(ctx, user.Id)
	// 			if err != nil {
	// 				s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
	// 			}
	// 			ret.Groups[i].Users = append(ret.Groups[i].Users, eventModel.UserStats{
	// 				UserId:   user.Id,
	// 				Username: userInfo.Email,
	// 				Avatar:   userInfo.Avatar,
	// 				Points:   userStats.Points,
	// 			})
	// 		}
	// 	}
	// }

	return nil, nil
}
