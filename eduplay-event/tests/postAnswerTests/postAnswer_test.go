package tests

import (
	"testing"

	dto "eduplay-event/internal/generated"
	event "eduplay-event/internal/pkg/usecase/events"

	"github.com/stretchr/testify/assert"
)

var (
	ansTask1 = &dto.Task{
		TaskId:      "1",
		Name:        "task1",
		Description: "description1",
		Type:        2,
		Options: []*dto.TaskOption{
			{
				OptionId:  "1",
				Value:     "option1",
				IsCorrect: true,
			},
			{
				OptionId:  "2",
				Value:     "option2",
				IsCorrect: true,
			},
			{
				OptionId:  "3",
				Value:     "option3",
				IsCorrect: false,
			},
		},
		Files:         nil,
		Points:        10,
		Time:          0,
		PartialPoints: false,
		BlockId:       "block1Id",
		Order:         1,
	}

	corrAnswers1 = []string{"option1", "option2"}

	corrAnswersIds1 = []string{"1", "2"}

	allOptions1 = map[string]string{
		"1": "option1",
		"2": "option2",
		"3": "option3",
	}
)

func TestMultitaskPoints_AllCorrAnswers(t *testing.T) {

	inAns := &dto.Answer{
		TaskId: ansTask1.TaskId,
		UserId: "user1",
		Answer: []string{"1", "2"},
		Points: 0,
		Status: "",
	}

	ans := &dto.Answer{
		TaskId:        ansTask1.TaskId,
		UserId:        "user1",
		Answer:        make([]string, 0),
		AnswerIds:     []string{"1", "2"},
		Points:        0,
		Status:        "",
		RightAnswer:   corrAnswers1,
		RightAnswerId: corrAnswersIds1,
	}

	ans = event.CountMulTaskPoints(inAns, ansTask1, corrAnswers1, corrAnswersIds1, allOptions1, ans)

	assert.Equal(t, "correct", ans.Status)
	assert.Equal(t, 10, int(ans.Points))
}

func TestMultitaskPoints_PartialCorrAnswer(t *testing.T) {

	inAns := &dto.Answer{
		TaskId: ansTask1.TaskId,
		UserId: "user1",
		Answer: []string{"1"},
		Points: 0,
		Status: "",
	}

	ans := &dto.Answer{
		TaskId:        ansTask1.TaskId,
		UserId:        "user1",
		Answer:        make([]string, 0),
		AnswerIds:     []string{"1"},
		Points:        0,
		Status:        "",
		RightAnswer:   corrAnswers1,
		RightAnswerId: corrAnswersIds1,
	}

	ans = event.CountMulTaskPoints(inAns, ansTask1, corrAnswers1, corrAnswersIds1, allOptions1, ans)

	assert.Equal(t, "partial", ans.Status)
	assert.Equal(t, 5, int(ans.Points))
}

func TestMultitaskPoints_TooManyOptions(t *testing.T) {

	inAns := &dto.Answer{
		TaskId: ansTask1.TaskId,
		UserId: "user1",
		Answer: []string{"1", "2", "3"},
		Points: 0,
		Status: "",
	}

	ans := &dto.Answer{
		TaskId:        ansTask1.TaskId,
		UserId:        "user1",
		Answer:        make([]string, 0),
		AnswerIds:     []string{"1", "2", "3"},
		Points:        0,
		Status:        "",
		RightAnswer:   corrAnswers1,
		RightAnswerId: corrAnswersIds1,
	}

	ans = event.CountMulTaskPoints(inAns, ansTask1, corrAnswers1, corrAnswersIds1, allOptions1, ans)

	assert.Equal(t, "partial", ans.Status)
	assert.Equal(t, 5, int(ans.Points))
}

func TestMultitaskPoints_NoCorrAnswers(t *testing.T) {

	inAns := &dto.Answer{
		TaskId: ansTask1.TaskId,
		UserId: "user1",
		Answer: []string{"3"},
		Points: 0,
		Status: "",
	}

	ans := &dto.Answer{
		TaskId:        ansTask1.TaskId,
		UserId:        "user1",
		Answer:        make([]string, 0),
		AnswerIds:     []string{"3"},
		Points:        0,
		Status:        "",
		RightAnswer:   corrAnswers1,
		RightAnswerId: corrAnswersIds1,
	}

	ans = event.CountMulTaskPoints(inAns, ansTask1, corrAnswers1, corrAnswersIds1, allOptions1, ans)

	assert.Equal(t, "incorrect", ans.Status)
	assert.Equal(t, 0, int(ans.Points))
}

func TestMultitaskPoints_ParticalCorrect_Excess(t *testing.T) {

	inAns := &dto.Answer{
		TaskId: ansTask1.TaskId,
		UserId: "user1",
		Answer: []string{"2", "3"},
		Points: 0,
		Status: "",
	}

	ans := &dto.Answer{
		TaskId:        ansTask1.TaskId,
		UserId:        "user1",
		Answer:        make([]string, 0),
		AnswerIds:     []string{"2", "3"},
		Points:        0,
		Status:        "",
		RightAnswer:   corrAnswers1,
		RightAnswerId: corrAnswersIds1,
	}

	ans = event.CountMulTaskPoints(inAns, ansTask1, corrAnswers1, corrAnswersIds1, allOptions1, ans)

	assert.Equal(t, "partial", ans.Status)
	assert.Equal(t, 0, int(ans.Points))
}
