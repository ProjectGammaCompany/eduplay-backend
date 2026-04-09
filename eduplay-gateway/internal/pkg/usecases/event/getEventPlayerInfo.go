package event

import (
	"context"
	eventDto "eduplay-gateway/internal/generated/clients/event"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"log/slog"
)

func (s *UseCase) GetEventPlayerInfo(ctx context.Context, userId string, eventId string) (*eventModel.EventPlayerInfo, error) {
	const op = "event.UseCase.GetEventPlayerInfo"

	s.log.With(slog.String("op", op)).Info("attempting to get event player info")

	playerInfo := &eventModel.EventPlayerInfo{}

	event, err := s.eventClient.GetEvent(ctx, &eventDto.Id{Id: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event", slog.String("error", err.Error()))
		return nil, err
	}

	collaborators, err := s.eventClient.GetCollaborators(ctx, &eventDto.Id{Id: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get collaborators", slog.String("error", err.Error()))
		return nil, err
	}

	eventForUser, err := s.eventClient.GetEventForUser(ctx, &eventDto.UserEventIds{UserId: userId, EventId: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event for user", slog.String("error", err.Error()))
		return nil, err
	}

	ownerProfile, err := s.userClient.GetProfile(ctx, event.OwnerId)
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get owner profile", slog.String("error", err.Error()))
		return nil, err
	}

	userStatus, err := s.eventClient.GetUserStatus(ctx, &eventDto.UserEventIds{UserId: userId, EventId: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get user status", slog.String("error", err.Error()))
		return nil, err
	}

	if event.GroupEvent {
		userGroup, err := s.eventClient.GetUserGroup(ctx, &eventDto.UserEventIds{UserId: userId, EventId: eventId})
		if err != nil {
			s.log.With(slog.String("op", op)).Error("failed to get user group", slog.String("error", err.Error()))
			return nil, err
		}

		if userGroup != nil || userGroup.GroupId != "" {
			playerInfo.NeedGroup = true
		}
	}

	rated := true
	rate, err := s.eventClient.GetEventUserRating(ctx, &eventDto.UserEventIds{UserId: userId, EventId: eventId})
	if err != nil {
		s.log.With(slog.String("op", op)).Error("failed to get event user rating", slog.String("error", err.Error()))
		return nil, err
	}
	if rate.Message == "-1" {
		rated = false
	}

	playerInfo.EventId = event.EventId
	playerInfo.Title = event.Title
	playerInfo.Description = event.Description
	playerInfo.Cover = event.Cover
	playerInfo.Tags = eventModel.TagsFromDto(eventForUser.Tags).Tags
	playerInfo.Authors = eventModel.CollaboratorsFromDto(collaborators)
	playerInfo.Rate = eventForUser.Rate
	playerInfo.Favorite = eventForUser.Favorite
	playerInfo.Status = userStatus.Message
	playerInfo.Authors = append(playerInfo.Authors, eventModel.Collaborator{Id: event.OwnerId, Email: ownerProfile.Email, Avatar: ownerProfile.Avatar})
	playerInfo.CanBeDownloaded = event.AllowDownloading
	playerInfo.Status = userStatus.Message
	playerInfo.IsPrivate = event.Private
	playerInfo.Rated = rated

	playerInfo.LastEditionDate = event.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000")
	if playerInfo.LastEditionDate == "01.01.1970 00:00:00.000" {
		playerInfo.LastEditionDate = ""
	}
	playerInfo.StartDate = event.StartDate.AsTime().Format("02.01.2006 15:04:05.000")
	if playerInfo.StartDate == "01.01.1970 00:00:00.000" {
		playerInfo.StartDate = ""
	}
	playerInfo.EndDate = event.EndDate.AsTime().Format("02.01.2006 15:04:05.000")
	if playerInfo.EndDate == "01.01.1970 00:00:00.000" {
		playerInfo.EndDate = ""
	}

	s.log.With(slog.String("op", op)).Info("got event player info", slog.Any("event", event.EventId))

	return playerInfo, nil
}
