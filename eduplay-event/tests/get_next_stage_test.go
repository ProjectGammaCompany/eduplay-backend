package tests

import (
	"context"
	"io"
	"log/slog"
	"testing"

	dto "eduplay-event/internal/generated"
	event "eduplay-event/internal/pkg/usecase/events"
	"eduplay-event/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var linkId = "linkId"
var currTaskId = ""
var currBlockId = ""

// var startTime *timestamppb.Timestamp
// var finished = false

var event1 = &dto.PostEventIn{
	EventId:    "eventId",
	GroupEvent: true,
}

var block1Conditions = []*dto.Condition{
	{
		ConditionId:     "condition1Id",
		PreviousBlockId: "block1Id",
		NextBlockId:     "block2Id",
		NextBlockOrder:  2,
		GroupIds:        []string{"group1Id", "group2Id"},
		Min:             0,
		Max:             5,
	},
	{
		ConditionId:     "condition2Id",
		PreviousBlockId: "block1Id",
		NextBlockId:     "block3Id",
		NextBlockOrder:  3,
		GroupIds:        []string{"group1Id", "group2Id"},
		Min:             6,
		Max:             -1,
	},
}

var block1 = &dto.PostEventBlockIn{
	BlockId:       "block1Id",
	EventId:       "eventId",
	Name:          "block1",
	Order:         1,
	IsParallel:    false,
	ShowPoints:    false,
	ShowAnswers:   false,
	PartialPoints: false,
}

var block2 = &dto.PostEventBlockIn{
	BlockId:       "block2Id",
	EventId:       "eventId",
	Name:          "block2",
	Order:         2,
	IsParallel:    true,
	ShowPoints:    false,
	ShowAnswers:   false,
	PartialPoints: false,
}

var block3 = &dto.PostEventBlockIn{
	BlockId:       "block3Id",
	EventId:       "eventId",
	Name:          "block3",
	Order:         3,
	IsParallel:    false,
	ShowPoints:    false,
	ShowAnswers:   false,
	PartialPoints: false,
}

var eventBlocks = &dto.GetEventBlocksOut{
	Blocks: []*dto.BlockInfo{
		{
			BlockId:    "block1Id",
			Name:       "block1",
			Order:      1,
			IsParallel: false,
		},
		{
			BlockId:    "block2Id",
			Name:       "block2",
			Order:      2,
			IsParallel: true,
		},
		{
			BlockId:    "block3Id",
			Name:       "block3",
			Order:      3,
			IsParallel: false,
		},
	},
}

var block1TasksShort = []*dto.NextStageTaskShort{
	{
		TaskId:      "task1Id",
		Name:        "task1",
		Time:        0,
		IsCompleted: false,
		Order:       1,
		Description: "description1",
		Type:        1,
	},
	{
		TaskId:      "task2Id",
		Name:        "task2",
		Time:        0,
		IsCompleted: false,
		Order:       2,
		Description: "description2",
		Type:        1,
	},
}

var block2TasksShort = []*dto.NextStageTaskShort{
	{
		TaskId:      "task3Id",
		Name:        "task3",
		Time:        0,
		IsCompleted: false,
		Order:       1,
		Description: "description3",
		Type:        0,
	},
}

var block3TasksShort = []*dto.NextStageTaskShort{}

var task1 = &dto.Task{
	TaskId:        "task1Id",
	Name:          "task1",
	Description:   "description1",
	Type:          1,
	Options:       nil,
	Files:         nil,
	Points:        6,
	Time:          0,
	PartialPoints: false,
	BlockId:       "block1Id",
	Order:         1,
}

var task2 = &dto.Task{
	TaskId:        "task2Id",
	Name:          "task2",
	Description:   "description2",
	Type:          1,
	Options:       nil,
	Files:         nil,
	Points:        6,
	Time:          0,
	PartialPoints: false,
	BlockId:       "block1Id",
	Order:         2,
}

// var task3 = &dto.Task{
// 	TaskId:        "task3Id",
// 	Name:          "task3",
// 	Description:   "description3",
// 	Type:          0,
// 	Options:       nil,
// 	Files:         nil,
// 	Points:        0,
// 	Time:          0,
// 	PartialPoints: false,
// 	BlockId:       "block2Id",
// 	Order:         1,
// }

var answer1 = &dto.Answer{
	UserId:    "user1",
	TaskId:    "task1Id",
	Answer:    []string{"answer1", "answer2"},
	AnswerIds: nil,
	Points:    6,
	Status:    "correct",
}

func TestGetNextStage_Finished(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)

	mockStorage.On("GetNextStage", ctx, mock.Anything).
		Return(linkId, currTaskId, currBlockId, true, timestamppb.Now(), nil)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := event.New(logger, mockStorage, "secret")

	res, err := uc.GetNextStage(ctx, &dto.UserEventIds{
		UserId:  "user1",
		EventId: "event1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "end", res.Type)

	mockStorage.AssertExpectations(t)
}

func TestGetNextStage_NewParticipant(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockStorage.
		On("GetNextStage", ctx, mock.Anything).
		Return(linkId, currTaskId, currBlockId, false, nil, nil)

	mockStorage.
		On("PutNextStage", ctx, mock.Anything).
		Return("updated", nil)

	mockStorage.
		On("GetEventBlocks", ctx, "eventId").
		Return(eventBlocks, nil)

	mockStorage.
		On("GetBlockInfo", ctx, "block1Id").
		Return(block1, nil)

	mockStorage.
		On("GetUserBlockTasksShort", ctx, "block1Id", "userId").
		Return(block1TasksShort, nil)

	mockStorage.
		On("GetTaskById", ctx, "task1Id").
		Return(task1, nil)

	uc := event.New(logger, mockStorage, "secret")

	nextStage, err := uc.GetNextStage(ctx, &dto.UserEventIds{
		EventId: "eventId",
		UserId:  "userId",
	})

	assert.NoError(t, err)
	assert.NotNil(t, nextStage)
	assert.NotEmpty(t, nextStage.Task)
	assert.Equal(t, "task", nextStage.Type)
	assert.Equal(t, task1.TaskId, nextStage.Task.TaskId)

	mockStorage.AssertExpectations(t)
}

func TestGetNextStage_SameTask(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockStorage.
		On("GetNextStage", ctx, mock.Anything).
		Return(linkId, block1TasksShort[0].TaskId, eventBlocks.Blocks[0].BlockId, false, timestamppb.Now(), nil)

	mockStorage.
		On("GetTaskAnswer", ctx, "task1Id", "userId").
		Return(nil, nil)

	mockStorage.
		On("GetTaskById", ctx, "task1Id").
		Return(task1, nil)

	uc := event.New(logger, mockStorage, "secret")

	nextStage, err := uc.GetNextStage(ctx, &dto.UserEventIds{
		EventId: "eventId",
		UserId:  "userId",
	})

	assert.NoError(t, err)
	assert.NotNil(t, nextStage)
	assert.NotEmpty(t, nextStage.Task)
	assert.Equal(t, "task", nextStage.Type)
	assert.Equal(t, task1.TaskId, nextStage.Task.TaskId)

	mockStorage.AssertExpectations(t)
}

func TestGetNextStage_NextTask(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockStorage.
		On("GetNextStage", ctx, mock.Anything).
		Return(linkId, block1TasksShort[0].TaskId, eventBlocks.Blocks[0].BlockId, false, nil, nil)

	mockStorage.
		On("GetTaskAnswer", ctx, "task1Id", "userId").
		Return(answer1, nil)

	mockStorage.
		On("GetTaskById", ctx, "task1Id").
		Return(task1, nil)

	mockStorage.
		On("GetTaskById", ctx, "task2Id").
		Return(task2, nil)

	mockStorage.
		On("GetBlockInfo", ctx, "block1Id").
		Return(block1, nil)

	mockStorage.
		On("GetUserBlockTasksShort", ctx, "block1Id", "userId").
		Return(block1TasksShort, nil)

	mockStorage.
		On("PutNextStage", ctx, mock.Anything).
		Return("updated", nil)

	uc := event.New(logger, mockStorage, "secret")

	nextStage, err := uc.GetNextStage(ctx, &dto.UserEventIds{
		EventId: "eventId",
		UserId:  "userId",
	})

	assert.NoError(t, err)
	assert.NotNil(t, nextStage)
	assert.NotEmpty(t, nextStage.Task)
	assert.Equal(t, "task", nextStage.Type)
	assert.Equal(t, task2.TaskId, nextStage.Task.TaskId)

	mockStorage.AssertExpectations(t)
}

func TestGetNextStage_NextBlockParallel(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockStorage.
		On("GetNextStage", ctx, mock.Anything).
		Return(linkId, block1TasksShort[len(block1TasksShort)-1].TaskId, eventBlocks.Blocks[0].BlockId, false, nil, nil)

	mockStorage.
		On("GetTaskAnswer", ctx, "task2Id", "userId").
		Return(answer1, nil)

	mockStorage.
		On("GetBlockInfo", ctx, "block1Id").
		Return(block1, nil)

	mockStorage.
		On("GetBlockInfo", ctx, "block2Id").
		Return(block2, nil)

	mockStorage.
		On("GetTaskById", ctx, "task2Id").
		Return(task2, nil)

	mockStorage.
		On("GetUserBlockTasksShort", ctx, "block1Id", "userId").
		Return(block1TasksShort, nil)

	mockStorage.
		On("GetUserBlockTasksShort", ctx, "block2Id", "userId").
		Return(block2TasksShort, nil)

	mockStorage.
		On("GetEventBlocks", ctx, "eventId").
		Return(eventBlocks, nil)

	mockStorage.
		On("GetBlockConditionsFull", ctx, "block1Id").
		Return(&dto.BlockInfo{Conditions: block1Conditions}, nil)

	mockStorage.
		On("GetUserBlockPointsSum", ctx, "userId", "block1Id").
		Return(int64(5), nil)

	mockStorage.
		On("GetEvent", ctx, "eventId").
		Return(event1, nil)

	mockStorage.
		On("GetUserGroup", ctx, "userId", "eventId").
		Return(&dto.GetUserGroupOut{GroupId: "group1Id"}, nil)

	mockStorage.
		On("GetBlockMaxPoints", ctx, "block1Id").
		Return(int64(12), nil)

	mockStorage.
		On("PutNextStage", ctx, mock.Anything).
		Return("updated", nil)

	uc := event.New(logger, mockStorage, "secret")

	nextStage, err := uc.GetNextStage(ctx, &dto.UserEventIds{
		EventId: "eventId",
		UserId:  "userId",
	})

	assert.NoError(t, err)
	assert.NotNil(t, nextStage)
	assert.NotEmpty(t, nextStage.Block)
	assert.Equal(t, "block", nextStage.Type)
	assert.Equal(t, block2.BlockId, nextStage.Block.BlockId)

	mockStorage.AssertExpectations(t)
}

func TestGetNextStage_NextBlockTask(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mocks.Storage)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockStorage.
		On("GetNextStage", ctx, mock.Anything).
		Return(linkId, block1TasksShort[len(block1TasksShort)-1].TaskId, eventBlocks.Blocks[0].BlockId, false, nil, nil)

	mockStorage.
		On("GetTaskAnswer", ctx, "task2Id", "userId").
		Return(answer1, nil)

	mockStorage.
		On("GetBlockInfo", ctx, "block1Id").
		Return(block1, nil)

	mockStorage.
		On("GetBlockInfo", ctx, "block3Id").
		Return(block3, nil)

	mockStorage.
		On("GetTaskById", ctx, "task2Id").
		Return(task2, nil)

	mockStorage.
		On("GetUserBlockTasksShort", ctx, "block1Id", "userId").
		Return(block1TasksShort, nil)

	mockStorage.
		On("GetUserBlockTasksShort", ctx, "block3Id", "userId").
		Return(block3TasksShort, nil)

	mockStorage.
		On("GetEventBlocks", ctx, "eventId").
		Return(eventBlocks, nil)

	mockStorage.
		On("GetBlockConditionsFull", ctx, "block1Id").
		Return(&dto.BlockInfo{Conditions: block1Conditions}, nil)

	mockStorage.
		On("GetUserBlockPointsSum", ctx, "userId", "block1Id").
		Return(int64(7), nil)

	mockStorage.
		On("GetEvent", ctx, "eventId").
		Return(event1, nil)

	mockStorage.
		On("GetUserGroup", ctx, "userId", "eventId").
		Return(&dto.GetUserGroupOut{GroupId: "group1Id"}, nil)

	mockStorage.
		On("GetBlockMaxPoints", ctx, "block1Id").
		Return(int64(12), nil)

	mockStorage.
		On("EndMe", ctx, "userId", "eventId").
		Return("updated to finished", nil)

	// mockStorage.
	// 	On("PutNextStage", ctx, mock.Anything).
	// 	Return("updated", nil)

	uc := event.New(logger, mockStorage, "secret")

	nextStage, err := uc.GetNextStage(ctx, &dto.UserEventIds{
		EventId: "eventId",
		UserId:  "userId",
	})

	assert.NoError(t, err)
	assert.NotNil(t, nextStage)
	assert.Equal(t, "end", nextStage.Type)

	mockStorage.AssertExpectations(t)
}
