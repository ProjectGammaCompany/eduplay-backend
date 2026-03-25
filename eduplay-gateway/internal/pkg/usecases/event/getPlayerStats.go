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
		if event.GroupEvent {
			userGroup, err := s.eventClient.GetUserGroup(ctx, &eventDto.UserEventIds{UserId: in.UserId, EventId: in.EventId})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get user group", slog.String("error", err.Error()))
				return nil, err
			}

			groupUsers, err := s.eventClient.GetGroupUsers(ctx, &eventDto.Id{Id: userGroup.GroupId})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get group users", slog.String("error", err.Error()))
				return nil, err
			}

			for i, user := range groupUsers.Users {
				if user.Id == in.UserId {
					groupUsers.Users[i].Current = true
				}

				userStats, err := s.eventClient.GetUserStats(ctx, &eventDto.UserEventIds{UserId: user.Id, EventId: in.EventId})
				if err != nil {
					s.log.With(slog.String("op", op)).Error("failed to get player stats", slog.String("error", err.Error()))
					return nil, err
				}

				userInfo, err := s.userClient.GetProfile(ctx, userStats.Id)
				if err != nil {
					s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
					return nil, err
				}

				groupUsers.Users[i].Email = userInfo.Email
				groupUsers.Users[i].Avatar = userInfo.Avatar
				groupUsers.Users[i].Points = userStats.Points
			}

			groupStats := make([]eventModel.GroupStats, 0)
			groupStats = append(groupStats, eventModel.GroupStatsFromDto(groupUsers))

			return &eventModel.PlayerStats{
				FullStats:  false,
				GroupEvent: true,
				Users:      nil,
				Groups:     groupStats,
			}, nil
		}

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

	if event.GroupEvent {
		eventGroups, err := s.eventClient.GetGroups(ctx, &eventDto.Id{Id: in.EventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get event groups", slog.String("error", err.Error()))
			return nil, err
		}

		groupStats := make([]eventModel.GroupStats, 0)

		for _, group := range eventGroups.Groups {
			groupStat := eventModel.GroupStats{
				GroupId: group.Id,
				Name:    group.Login,
				Users:   make([]eventModel.UserStats, 0),
			}

			groupUsers, err := s.eventClient.GetGroupUsers(ctx, &eventDto.Id{Id: group.Id})
			if err != nil {
				s.log.With(slog.String("op", op)).Error("failed to get group users", slog.String("error", err.Error()))
				return nil, err
			}

			for _, user := range groupUsers.Users {
				userStat := eventModel.UserStats{
					UserId: user.Id,
				}

				if user.Id == in.UserId {
					userStat.Current = true
				}

				userStats, err := s.eventClient.GetUserStats(ctx, &eventDto.UserEventIds{UserId: user.Id, EventId: in.EventId})
				if err != nil {
					s.log.With(slog.String("op", op)).Error("failed to get player stats", slog.String("error", err.Error()))
					return nil, err
				}

				userInfo, err := s.userClient.GetProfile(ctx, userStats.Id)
				if err != nil {
					s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
					return nil, err
				}

				userStat.Username = userInfo.Email
				userStat.Avatar = userInfo.Avatar
				userStat.Points = userStats.Points

				groupStat.Users = append(groupStat.Users, userStat)
			}

			groupStats = append(groupStats, groupStat)
		}

		return &eventModel.PlayerStats{
			FullStats:  true,
			GroupEvent: true,
			Users:      nil,
			Groups:     groupStats,
		}, nil
	}

	eventUsers, err := s.eventClient.GetEventUsers(ctx, &eventDto.Id{Id: in.EventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event users", slog.String("error", err.Error()))
		return nil, err
	}

	userStats := make([]eventModel.UserStats, 0)

	for _, user := range eventUsers.Users {
		userStat := eventModel.UserStats{
			UserId: user.Id,
		}

		if user.Id == in.UserId {
			userStat.Current = true
		}

		userStatDto, err := s.eventClient.GetUserStats(ctx, &eventDto.UserEventIds{UserId: user.Id, EventId: in.EventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get player stats", slog.String("error", err.Error()))
			return nil, err
		}
		userStat.Points = userStatDto.Points

		userInfo, err := s.userClient.GetProfile(ctx, userStatDto.Id)
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user profile", slog.String("error", err.Error()))
			return nil, err
		}

		userStat.Username = userInfo.Email
		userStat.Avatar = userInfo.Avatar

		userStats = append(userStats, userStat)
	}

	return &eventModel.PlayerStats{
		FullStats:  true,
		GroupEvent: false,
		Users:      userStats,
		Groups:     nil,
	}, nil
}
