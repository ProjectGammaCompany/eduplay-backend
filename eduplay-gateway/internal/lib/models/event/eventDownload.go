package eventModel

type EventDownloadFull struct {
	EventDownload          EventDownload           `json:"event"`
	BlocksDownload         []BlockDownload         `json:"blocks"`
	ConditionsDownload     []ConditionDownload     `json:"conditions"`
	GroupsDownload         []GroupDownload         `json:"groups"`
	TasksDownload          []TaskDownload          `json:"tasks"`
	OptionsDownload        []OptionDownload        `json:"options"`
	CorrectAnswersDownload []CorrectAnswerDownload `json:"correctAnswers"`
	Files                  []string                `json:"files"`
}

type EventDownload struct {
	EventId         string   `json:"eventId"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	Cover           string   `json:"cover"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	LastEditionDate string   `json:"lastEditionDate"`
	GroupEvent      bool     `json:"groupEvent"`
	AuthorId        []string `json:"authorId"`
}

type BlockDownload struct {
	BlockId       string `json:"blockId"`
	Name          string `json:"name"`
	BlockOrder    int64  `json:"blockOrder"`
	IsParallel    bool   `json:"isParallel"`
	ShowPoints    bool   `json:"showPoints"`
	ShowAnswers   bool   `json:"showAnswers"`
	PartialPoints bool   `json:"PartialPoints"`
	EventId       string `json:"eventId"`
}

type ConditionDownload struct {
	ConditionId string   `json:"conditionId"`
	PrevBlockId string   `json:"prevBlockId"`
	NextBlockId string   `json:"nextBlockId"`
	GroupName   []string `json:"groupName"`
	Min         *int64   `json:"min"`
	Max         *int64   `json:"max"`
}

type GroupDownload struct {
	GroupId  string `json:"groupId"`
	EventId  string `json:"eventId"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TaskDownload struct {
	TaskId        string   `json:"taskId"`
	BlockId       string   `json:"blockId"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	TaskType      int64    `json:"type"`
	Files         []string `json:"files"`
	Points        int64    `json:"points"`
	PartialPoints bool     `json:"partialPoints"`
	Time          int64    `json:"time"`
	Order         int64    `json:"taskOrder"`
}

type OptionDownload struct {
	OptionId string `json:"optionId"`
	TaskId   string `json:"taskId"`
	Value    string `json:"value"`
}

type CorrectAnswerDownload struct {
	TaskId string   `json:"taskId"`
	Values []string `json:"values"`
}
