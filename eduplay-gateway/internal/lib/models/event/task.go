package eventModel

import (
	dto "eduplay-gateway/internal/generated/clients/event"
)

type TaskOption struct {
	OptionId  string `json:"id"`
	Value     string `json:"value"`
	IsCorrect bool   `json:"isCorrect"`
}

type Task struct {
	TaskId        string        `json:"id"`
	BlockId       string        `json:"blockId"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	TaskType      int64         `json:"type"`
	Options       []*TaskOption `json:"options"`
	Files         []string      `json:"files"`
	Points        int64         `json:"points"`
	Time          int64         `json:"time"`
	PartialPoints bool          `json:"partialPoints"`
}

func TaskToDto(task *Task) *dto.Task {
	files := make([]string, 0)

	if task.Files != nil || len(task.Files) > 0 {
		files = task.Files
	}

	return &dto.Task{
		TaskId:        task.TaskId,
		BlockId:       task.BlockId,
		Name:          task.Name,
		Description:   task.Description,
		Type:          task.TaskType,
		Options:       TaskOptionsToDto(task.Options),
		Files:         files,
		Points:        task.Points,
		Time:          task.Time,
		PartialPoints: task.PartialPoints,
	}
}

func TaskFromDto(task *dto.Task) *Task {
	files := make([]string, 0)

	if task.Files != nil || len(task.Files) > 0 {
		files = task.Files
	}

	return &Task{
		TaskId:        task.TaskId,
		BlockId:       task.BlockId,
		Name:          task.Name,
		Description:   task.Description,
		TaskType:      task.Type,
		Options:       TaskOptionsFromDto(task.Options),
		Files:         files,
		Points:        task.Points,
		Time:          task.Time,
		PartialPoints: task.PartialPoints,
	}
}

func TaskOptionsFromDto(taskOptions []*dto.TaskOption) []*TaskOption {
	options := make([]*TaskOption, 0)
	for _, taskOption := range taskOptions {
		options = append(options, &TaskOption{
			OptionId:  taskOption.OptionId,
			Value:     taskOption.Value,
			IsCorrect: taskOption.IsCorrect,
		})
	}
	return options
}

func TaskOptionsToDto(taskOptions []*TaskOption) []*dto.TaskOption {
	dtoTaskOptions := make([]*dto.TaskOption, 0)
	for _, taskOption := range taskOptions {
		dtoTaskOptions = append(dtoTaskOptions, &dto.TaskOption{
			OptionId:  taskOption.OptionId,
			Value:     taskOption.Value,
			IsCorrect: taskOption.IsCorrect,
		})
	}
	return dtoTaskOptions
}

type ShortTask struct {
	TaskId string `json:"id"`
	Name   string `json:"name"`
	Order  int64  `json:"order"`
}

type BlockTasksList struct {
	Tasks []ShortTask `json:"tasks"`
}

func BlockTasksListFromDto(blockTasksList *dto.Tasks) *BlockTasksList {
	tasks := make([]ShortTask, 0)
	for _, task := range blockTasksList.Tasks {
		tasks = append(tasks, ShortTask{
			TaskId: task.TaskId,
			Name:   task.Name,
			Order:  task.Order,
		})
	}
	return &BlockTasksList{Tasks: tasks}
}

// message Answer {
//     string taskId = 1;
//     string userId = 2;
//     repeated string answer = 3;
//     int64 points = 4;
//     string status = 5;
//     repeated string rightAnswer = 6;
// }

type Answer struct {
	TaskId      string   `json:"taskId"`
	UserId      string   `json:"userId"`
	Answer      []string `json:"answer"`
	Points      int64    `json:"points"`
	Status      string   `json:"status"`
	RightAnswer []string `json:"rightAnswer"`
}

func AnswerFromDto(answer *dto.Answer) *Answer {
	return &Answer{
		TaskId:      answer.TaskId,
		UserId:      answer.UserId,
		Answer:      answer.Answer,
		Points:      answer.Points,
		Status:      answer.Status,
		RightAnswer: answer.RightAnswer,
	}
}

func AnswerToDto(answer *Answer) *dto.Answer {
	return &dto.Answer{
		TaskId:      answer.TaskId,
		UserId:      answer.UserId,
		Answer:      answer.Answer,
		Points:      answer.Points,
		Status:      answer.Status,
		RightAnswer: answer.RightAnswer,
	}
}
