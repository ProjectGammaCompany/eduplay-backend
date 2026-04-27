package usecase_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	eventDto "eduplay-gateway/internal/generated/clients/event"
	userDto "eduplay-gateway/internal/generated/clients/user"
	event "eduplay-gateway/internal/pkg/usecases/event"
	"eduplay-gateway/tests/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	event1 = &eventDto.PostEventIn{
		EventId:          "event1Id",
		Title:            "event1",
		Description:      "description1",
		Tags:             []string{"tag1Id", "tag2Id"},
		Cover:            "cover1",
		StartDate:        timestamppb.Now(),
		EndDate:          timestamppb.Now(),
		Private:          true,
		Password:         "password1",
		OwnerId:          "owner1Id",
		LastEditionDate:  timestamppb.Now(),
		AllowDownloading: true,
		GroupEvent:       true,
		Rating:           true,
		EventRating:      0,
	}

	event2 = &eventDto.PostEventIn{
		EventId:          "event2Id",
		Title:            "event2",
		Description:      "description2",
		Tags:             []string{"tag1Id", "tag2Id"},
		Cover:            "cover2",
		StartDate:        timestamppb.Now(),
		EndDate:          timestamppb.Now(),
		Private:          true,
		Password:         "password2",
		OwnerId:          "owner2Id",
		LastEditionDate:  timestamppb.Now(),
		AllowDownloading: true,
		GroupEvent:       false,
		Rating:           true,
		EventRating:      0,
	}

	collaborators = &eventDto.GetCollaboratorsOut{
		Users: []*eventDto.User{
			{
				Id:     "user1Id",
				Email:  "user1@email",
				Avatar: "avatar1",
			},
			{
				Id:     "user2Id",
				Email:  "user2@email",
				Avatar: "avatar2",
			},
		},
	}

	eventForUser1 = &eventDto.GetPublicEvent{
		Tags: []*eventDto.Tag{
			{
				Id:   "tag1Id",
				Name: "tag1",
			},
			{
				Id:   "tag2Id",
				Name: "tag2",
			},
		},
		Rate:     0,
		Favorite: false,
	}

	ownerProfile = &userDto.Profile{
		UserId:   "owner1Id",
		Email:    "owner1@email",
		Avatar:   "avatar1",
		UserName: "owner1",
	}

	user1Status = &eventDto.UserStatus{
		UserId:  "user1Id",
		EventId: "event1Id",
		Status:  "in progress",
	}

	user1GroupOut = &eventDto.GetUserGroupOut{
		GroupId: "group1Id",
		Name:    "group1",
	}
)

func TestGetEventSettings_UserInGroup(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event1Id"}).
		Return(event1, nil)

	mockEventClient.On("GetCollaborators", ctx, &eventDto.Id{Id: "event1Id"}).
		Return(collaborators, nil)

	mockEventClient.On("GetEventForUser", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(eventForUser1, nil)

	mockUserClient.On("GetProfile", ctx, "owner1Id").
		Return(ownerProfile, nil)

	mockEventClient.
		On("GetUserStatus", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(user1Status, nil)

	mockEventClient.On("GetUserGroup", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(user1GroupOut, nil)

	mockEventClient.On("GetEventUserRating", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(&eventDto.MessageOut{Message: "-1"}, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetEventPlayerInfo(ctx, "user1Id", "event1Id")

	assert.NoError(t, err)
	assert.Equal(t, false, res.NeedGroup)
	assert.Equal(t, false, res.Rated)

	mockEventClient.AssertExpectations(t)
}

func TestGetEventSettings_UserNotInGroup(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event1Id"}).
		Return(event1, nil)

	mockEventClient.On("GetCollaborators", ctx, &eventDto.Id{Id: "event1Id"}).
		Return(collaborators, nil)

	mockEventClient.On("GetEventForUser", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(eventForUser1, nil)

	mockUserClient.On("GetProfile", ctx, "owner1Id").
		Return(ownerProfile, nil)

	mockEventClient.
		On("GetUserStatus", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(user1Status, nil)

	mockEventClient.On("GetUserGroup", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(nil, nil)

	mockEventClient.On("GetEventUserRating", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event1Id"}).
		Return(&eventDto.MessageOut{Message: "-1"}, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetEventPlayerInfo(ctx, "user1Id", "event1Id")

	assert.NoError(t, err)
	assert.Equal(t, true, res.NeedGroup)
	assert.Equal(t, false, res.Rated)

	mockEventClient.AssertExpectations(t)
}

func TestGetEventSettings_NotGroupEvent(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event2Id"}).
		Return(event2, nil)

	mockEventClient.On("GetCollaborators", ctx, &eventDto.Id{Id: "event2Id"}).
		Return(collaborators, nil)

	mockEventClient.On("GetEventForUser", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event2Id"}).
		Return(eventForUser1, nil)

	mockUserClient.On("GetProfile", ctx, "owner2Id").
		Return(ownerProfile, nil)

	mockEventClient.
		On("GetUserStatus", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event2Id"}).
		Return(user1Status, nil)

	mockEventClient.On("GetEventUserRating", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event2Id"}).
		Return(&eventDto.MessageOut{Message: "-1"}, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetEventPlayerInfo(ctx, "user1Id", "event2Id")

	assert.NoError(t, err)
	assert.Equal(t, false, res.NeedGroup)
	assert.Equal(t, false, res.Rated)

	mockEventClient.AssertExpectations(t)
}

func TestGetEventSettings_Rated(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event2Id"}).
		Return(event2, nil)

	mockEventClient.On("GetCollaborators", ctx, &eventDto.Id{Id: "event2Id"}).
		Return(collaborators, nil)

	mockEventClient.On("GetEventForUser", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event2Id"}).
		Return(eventForUser1, nil)

	mockUserClient.On("GetProfile", ctx, "owner2Id").
		Return(ownerProfile, nil)

	mockEventClient.
		On("GetUserStatus", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event2Id"}).
		Return(user1Status, nil)

	mockEventClient.On("GetEventUserRating", ctx, &eventDto.UserEventIds{UserId: "user1Id", EventId: "event2Id"}).
		Return(&eventDto.MessageOut{Message: "4"}, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetEventPlayerInfo(ctx, "user1Id", "event2Id")

	assert.NoError(t, err)
	assert.Equal(t, false, res.NeedGroup)
	assert.Equal(t, true, res.Rated)

	mockEventClient.AssertExpectations(t)
}
