package usecase_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	eventDto "eduplay-gateway/internal/generated/clients/event"
	userDto "eduplay-gateway/internal/generated/clients/user"
	eventModel "eduplay-gateway/internal/lib/models/event"
	event "eduplay-gateway/internal/pkg/usecases/event"
	"eduplay-gateway/tests/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	event3 = &eventDto.PostEventIn{
		EventId:          "event3Id",
		Title:            "event3",
		Description:      "description3",
		Tags:             []string{"tag1Id", "tag2Id"},
		Cover:            "cover3",
		StartDate:        timestamppb.Now(),
		EndDate:          timestamppb.Now(),
		Private:          true,
		Password:         "password3",
		OwnerId:          "owner1Id",
		LastEditionDate:  timestamppb.Now(),
		AllowDownloading: true,
		GroupEvent:       true,
		Rating:           false,
		EventRating:      0,
	}

	event4 = &eventDto.PostEventIn{
		EventId:          "event4Id",
		Title:            "event4",
		Description:      "description4",
		Tags:             []string{"tag1Id", "tag2Id"},
		Cover:            "cover4",
		StartDate:        timestamppb.Now(),
		EndDate:          timestamppb.Now(),
		Private:          true,
		Password:         "password4",
		OwnerId:          "owner2Id",
		LastEditionDate:  timestamppb.Now(),
		AllowDownloading: true,
		GroupEvent:       false,
		Rating:           false,
		EventRating:      0,
	}

	event4User4Stats = &eventDto.User{
		Id:      "user4Id",
		Email:   "",
		Avatar:  "",
		Points:  10,
		Current: false,
	}

	event4User3Stats = &eventDto.User{
		Id:      "user3Id",
		Email:   "",
		Avatar:  "",
		Points:  11,
		Current: false,
	}

	event3User4Stats = &eventDto.User{
		Id:      "user4Id",
		Email:   "",
		Avatar:  "",
		Points:  15,
		Current: false,
	}

	event3User3Stats = &eventDto.User{
		Id:      "user3Id",
		Email:   "",
		Avatar:  "",
		Points:  20,
		Current: false,
	}

	user4 = &userDto.Profile{
		UserId:   "user4Id",
		Email:    "user4@email",
		Avatar:   "avatar4",
		UserName: "user4",
	}

	user3 = &userDto.Profile{
		UserId:   "user3Id",
		Email:    "user3@email",
		Avatar:   "avatar3",
		UserName: "user3",
	}

	user4Event3Group = &eventDto.GetUserGroupOut{
		GroupId: "group3Id",
		Name:    "group3",
	}

	groupUsers = &eventDto.GetGroupUsersOut{
		GroupId: "group3Id",
		Name:    "group3",
		Users: []*eventDto.User{
			{
				Id:      "user4Id",
				Email:   "",
				Avatar:  "",
				Points:  0,
				Current: false,
			},
			{
				Id:      "user3Id",
				Email:   "",
				Avatar:  "",
				Points:  0,
				Current: false,
			},
		},
	}

	event2Users = &eventDto.GetCollaboratorsOut{
		Users: []*eventDto.User{
			{
				Id:      "user4Id",
				Email:   "",
				Avatar:  "",
				Points:  0,
				Current: false,
			},
			{
				Id:      "user3Id",
				Email:   "",
				Avatar:  "",
				Points:  0,
				Current: false,
			},
		},
	}
)

func TestGetPlayerStats_FullStatsFalse_GroupFalse(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event4Id"}).
		Return(event4, nil)

	mockEventClient.On("GetUserStats", ctx, &eventDto.UserEventIds{UserId: "user4Id", EventId: "event4Id"}).
		Return(event4User4Stats, nil)

	mockUserClient.On("GetProfile", ctx, "user4Id").
		Return(user4, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetPlayerStats(ctx, &eventModel.UserEventIds{UserId: "user4Id", EventId: "event4Id"})

	assert.NoError(t, err)
	assert.Equal(t, false, res.FullStats)
	assert.Equal(t, false, res.GroupEvent)
	assert.Empty(t, res.Groups)
	assert.Equal(t, 1, len(res.Users))
	assert.Equal(t, "user4", res.Users[0].Username)
	assert.Equal(t, int64(10), res.Users[0].Points)

	mockEventClient.AssertExpectations(t)
}

func TestGetPlayerStats_FullStatsFalse_GroupTrue(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event3Id"}).
		Return(event3, nil)

	mockEventClient.On("GetUserStats", ctx, &eventDto.UserEventIds{UserId: "user4Id", EventId: "event3Id"}).
		Return(event3User4Stats, nil)

	mockEventClient.On("GetUserGroup", ctx, &eventDto.UserEventIds{UserId: "user4Id", EventId: "event3Id"}).
		Return(user4Event3Group, nil)

	mockUserClient.On("GetProfile", ctx, "user4Id").
		Return(user4, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetPlayerStats(ctx, &eventModel.UserEventIds{UserId: "user4Id", EventId: "event3Id"})

	assert.NoError(t, err)
	assert.Equal(t, false, res.FullStats)
	assert.Equal(t, true, res.GroupEvent)
	assert.Empty(t, res.Users)
	assert.Equal(t, 1, len(res.Groups))
	assert.Equal(t, "group3", res.Groups[0].Name)
	assert.Equal(t, 1, len(res.Groups[0].Users))
	assert.Equal(t, int64(15), res.Groups[0].Users[0].Points)

	mockEventClient.AssertExpectations(t)
}

func TestGetPlayerStats_FullStatsTrue_GroupFalse(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event2Id"}).
		Return(event2, nil)

	mockEventClient.On("GetEventUsers", ctx, &eventDto.Id{Id: "event2Id"}).
		Return(event2Users, nil)

	mockEventClient.On("GetUserStats", ctx, &eventDto.UserEventIds{UserId: "user4Id", EventId: "event2Id"}).
		Return(event4User4Stats, nil)

	mockEventClient.On("GetUserStats", ctx, &eventDto.UserEventIds{UserId: "user3Id", EventId: "event2Id"}).
		Return(event4User3Stats, nil)

	mockUserClient.On("GetProfile", ctx, "user4Id").
		Return(user4, nil)

	mockUserClient.On("GetProfile", ctx, "user3Id").
		Return(user3, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetPlayerStats(ctx, &eventModel.UserEventIds{UserId: "user4Id", EventId: "event2Id"})

	expectedUsers := []eventModel.UserStats{
		{
			UserId:   "user4Id",
			Username: "user4",
			Avatar:   "avatar4",
			Points:   10,
			Current:  true,
		},
		{
			UserId:   "user3Id",
			Username: "user3",
			Avatar:   "avatar3",
			Points:   11,
			Current:  false,
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, true, res.FullStats)
	assert.Equal(t, false, res.GroupEvent)
	assert.Empty(t, res.Groups)
	assert.Equal(t, 2, len(res.Users))
	assert.ElementsMatch(t, expectedUsers, res.Users)

	mockEventClient.AssertExpectations(t)
}

func TestGetPlayerStats_FullStatsTrue_GroupTrue(t *testing.T) {
	ctx := context.Background()

	mockEventClient := new(mocks.EventClient)
	mockUserClient := new(mocks.UserClient)

	mockEventClient.On("GetEvent", ctx, &eventDto.Id{Id: "event3Id"}).
		Return(event1, nil)

	mockEventClient.On("GetUserGroup", ctx, &eventDto.UserEventIds{UserId: "user4Id", EventId: "event3Id"}).
		Return(user4Event3Group, nil)

	mockEventClient.On("GetGroupUsers", ctx, &eventDto.Id{Id: "group3Id"}).
		Return(groupUsers, nil)

	mockEventClient.On("GetUserStats", ctx, &eventDto.UserEventIds{UserId: "user4Id", EventId: "event3Id"}).
		Return(event3User4Stats, nil)

	mockEventClient.On("GetUserStats", ctx, &eventDto.UserEventIds{UserId: "user3Id", EventId: "event3Id"}).
		Return(event3User3Stats, nil)

	mockUserClient.On("GetProfile", ctx, "user4Id").
		Return(user4, nil)

	mockUserClient.On("GetProfile", ctx, "user3Id").
		Return(user3, nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockEventClient, mockUserClient)

	res, err := uc.GetPlayerStats(ctx, &eventModel.UserEventIds{UserId: "user4Id", EventId: "event3Id"})

	expectedUsers := []eventModel.UserStats{
		{
			UserId:   "user4Id",
			Username: "user4",
			Avatar:   "avatar4",
			Points:   15,
			Current:  true,
		},
		{
			UserId:   "user3Id",
			Username: "user3",
			Avatar:   "avatar3",
			Points:   20,
			Current:  false,
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, true, res.FullStats)
	assert.Equal(t, true, res.GroupEvent)
	assert.Empty(t, res.Users)
	assert.Equal(t, 1, len(res.Groups))
	assert.Equal(t, 2, len(res.Groups[0].Users))
	assert.ElementsMatch(t, expectedUsers, res.Groups[0].Users)

	mockEventClient.AssertExpectations(t)
}
