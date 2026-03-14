package eventModel

import (
	dto "eduplay-gateway/internal/generated/clients/event"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PutEventIn struct {
	EventId     string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Tags            []Tag   `json:"tags"`
	Tags            []string `json:"tags"`
	Cover           string   `json:"cover"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	Private         bool     `json:"private"`
	Password        string   `json:"password"`
	GroupEvent      bool     `json:"groupEvent"`
	LastEditionDate string   `json:"lastEditionDate"`
	Groups          []Group  `json:"groups"`
	Rating          bool     `json:"rating"`
	// Collaborators []Collaborator `json:"collaborators"`
	Collaborators    []string `json:"collaborators"`
	AllowDownloading bool     `json:"allowDownloading"`
}

func PutEventInToDto(in *PutEventIn) (*dto.PutEventIn, error) {
	eventDto := &dto.PutEventIn{
		EventId:     in.EventId,
		Title:       in.Name,
		Description: in.Description,
		Tags:        in.Tags,
		Cover:       in.Cover,
		Private:     in.Private,
		Password:    in.Password,
		GroupEvent:  in.GroupEvent,
		Rating:      in.Rating,
		// Collaborators:   in.Collaborators,
		Collaborators:    in.Collaborators,
		AllowDownloading: in.AllowDownloading,
	}

	if in.StartDate != "" {
		startDate, err := time.Parse("02.01.2006 15:04", in.StartDate)
		if err != nil {
			return nil, err
		}
		eventDto.StartDate = timestamppb.New(startDate)
	}

	if in.EndDate != "" {
		endDate, err := time.Parse("02.01.2006 15:04", in.EndDate)
		if err != nil {
			return nil, err
		}
		eventDto.EndDate = timestamppb.New(endDate)
	}

	if in.LastEditionDate != "" {
		lastEditionDate, err := time.Parse("02.01.2006 15:04:05.000", in.LastEditionDate)
		if err != nil {
			return nil, err
		}
		eventDto.LastEditionDate = timestamppb.New(lastEditionDate)
	}

	gps := make([]*dto.Group, len(in.Groups))
	for i, group := range in.Groups {
		gps[i] = &dto.Group{
			Id:       group.Id,
			Login:    group.Login,
			Password: group.Password,
		}
	}

	eventDto.Groups = gps

	return eventDto, nil
}

type PostEventIn struct {
	EventId          string   `json:"id"`
	Title            string   `json:"title" validate:"required"`
	Description      string   `json:"description"`
	Tags             []string `json:"tags"`
	Cover            string   `json:"cover"`
	StartDate        string   `json:"startDate"`
	EndDate          string   `json:"endDate"`
	Private          bool     `json:"private"`
	Password         string   `json:"password"`
	OwnerId          string   `json:"ownerId"`
	LastEditionDate  string   `json:"lastEditionDate"`
	AllowDownloading bool     `json:"allowDownloading"`
	GroupEvent       bool     `json:"groupEvent"`
	Rating           bool     `json:"rating"`
}

func PostEventInToDto(in *PostEventIn) (*dto.PostEventIn, error) {
	eventDto := &dto.PostEventIn{
		EventId:          in.EventId,
		Title:            in.Title,
		Description:      in.Description,
		Tags:             in.Tags,
		Cover:            in.Cover,
		Private:          in.Private,
		Password:         in.Password,
		OwnerId:          in.OwnerId,
		AllowDownloading: in.AllowDownloading,
		GroupEvent:       in.GroupEvent,
		Rating:           in.Rating,
	}

	if in.StartDate != "" {
		startDate, err := time.Parse("02.01.2006 15:04", in.StartDate)
		if err != nil {
			return nil, err
		}
		eventDto.StartDate = timestamppb.New(startDate)
	}

	if in.EndDate != "" {
		endDate, err := time.Parse("02.01.2006 15:04", in.EndDate)
		if err != nil {
			return nil, err
		}
		eventDto.EndDate = timestamppb.New(endDate)
	}

	if in.LastEditionDate != "" {
		lastEditionDate, err := time.Parse("02.01.2006 15:04:05.000", in.LastEditionDate)
		if err != nil {
			return nil, err
		}
		eventDto.LastEditionDate = timestamppb.New(lastEditionDate)
	}

	return eventDto, nil
}

func PostEventInFromDto(in *dto.PostEventIn) *PostEventIn {
	event := &PostEventIn{
		EventId:          in.EventId,
		Title:            in.Title,
		Description:      in.Description,
		Tags:             in.Tags,
		Cover:            in.Cover,
		Private:          in.Private,
		Password:         in.Password,
		OwnerId:          in.OwnerId,
		StartDate:        in.StartDate.AsTime().Format("02.01.2006 15:04:05.000"),
		EndDate:          in.EndDate.AsTime().Format("02.01.2006 15:04:05.000"),
		LastEditionDate:  in.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000"),
		AllowDownloading: in.AllowDownloading,
		GroupEvent:       in.GroupEvent,
		Rating:           in.Rating,
	}

	if event.StartDate == "01.01.1970 00:00:00.000" {
		event.StartDate = ""
	}

	if event.EndDate == "01.01.1970 00:00:00.000" {
		event.EndDate = ""
	}

	if event.LastEditionDate == "01.01.1970 00:00:00.000" {
		event.LastEditionDate = ""
	}

	return event
}

type Id struct {
	Id string `json:"id"`
}

type PutFavorite struct {
	UserId   string `json:"userId"`
	EventId  string `json:"eventId" validate:"required"`
	Favorite bool   `json:"isFavorite"`
}

type Group struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GroupShort struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Groups struct {
	Groups []Group `json:"groups"`
}

type GroupsShort struct {
	Groups []GroupShort `json:"groups"`
}

func GroupsFromDto(in *dto.GetGroupsOut) *Groups {
	gps := make([]Group, len(in.Groups))
	for i, group := range in.Groups {
		gps[i] = Group{
			Id:       group.Id,
			Login:    group.Login,
			Password: group.Password,
		}
	}
	return &Groups{Groups: gps}
}

func GroupsShortFromDto(in *dto.GetGroupsOut) *GroupsShort {
	gps := make([]GroupShort, len(in.Groups))
	for i, group := range in.Groups {
		gps[i] = GroupShort{
			Id:   group.Id,
			Name: group.Login,
		}
	}
	return &GroupsShort{Groups: gps}
}

type PutGroupsIn struct {
	ConditionId string   `json:"conditionId"`
	GroupIds    []string `json:"groups"`
}

func PutGroupsInToDto(in *PutGroupsIn) *dto.PutListIn {
	return &dto.PutListIn{
		Id:   in.ConditionId,
		List: in.GroupIds,
	}
}

type Collaborator struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type GetEventSettings struct {
	EventId         string   `json:"id"`
	Title           string   `json:"title" validate:"required"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	Cover           string   `json:"cover"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	Private         bool     `json:"private"`
	Password        string   `json:"password"`
	LastEditionDate string   `json:"lastEditionDate"`
	Groups          []Group  `json:"groups"`
	Rating          bool     `json:"rating"`
	// Collaborators    []Collaborator `json:"collaborators"`
	Collaborators    []string `json:"collaborators"`
	AllowDownloading bool     `json:"allowDownloading"`
	OwnerId          string   `json:"ownerId"`
	GroupEvent       bool     `json:"groupEvent"`
}

func GetEventSettingsFromDto(event *dto.PostEventIn, groups *dto.GetGroupsOut, collaborators *dto.GetCollaboratorsOut) *GetEventSettings {
	gps := make([]Group, len(groups.Groups))
	for i, group := range groups.Groups {
		gps[i] = Group{
			Id:       group.Id,
			Login:    group.Login,
			Password: group.Password,
		}
	}

	collabs := CollaboratorsToString(collaborators)

	eventOut := &GetEventSettings{
		EventId:          event.EventId,
		Title:            event.Title,
		Description:      event.Description,
		Tags:             event.Tags,
		Cover:            event.Cover,
		StartDate:        event.StartDate.AsTime().Format("02.01.2006 15:04:05.000"),
		EndDate:          event.EndDate.AsTime().Format("02.01.2006 15:04:05.000"),
		Private:          event.Private,
		Password:         event.Password,
		LastEditionDate:  event.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000"),
		Groups:           gps,
		Rating:           event.Rating,
		Collaborators:    collabs,
		AllowDownloading: event.AllowDownloading,
		OwnerId:          event.OwnerId,
		GroupEvent:       event.GroupEvent,
	}

	if eventOut.StartDate == "01.01.1970 00:00:00.000" {
		eventOut.StartDate = ""
	}

	if eventOut.EndDate == "01.01.1970 00:00:00.000" {
		eventOut.EndDate = ""
	}

	if eventOut.LastEditionDate == "01.01.1970 00:00:00.000" {
		eventOut.LastEditionDate = ""
	}

	return eventOut
}

func CollaboratorsFromDto(collaborators *dto.GetCollaboratorsOut) []Collaborator {
	collabs := make([]Collaborator, len(collaborators.Users))
	for i, user := range collaborators.Users {
		collabs[i] = Collaborator{
			Id:     user.Id,
			Email:  user.Email,
			Avatar: user.Avatar,
		}
	}
	return collabs
}

func CollaboratorsToString(collaborators *dto.GetCollaboratorsOut) []string {
	collabs := make([]string, len(collaborators.Users))
	for i, user := range collaborators.Users {
		collabs[i] = user.Id
	}
	return collabs
}

type EventBaseFilters struct {
	Page            int64    `json:"page"`
	MaxOnPage       int64    `json:"maxOnPage"`
	Tags            []string `json:"tags"`
	DecliningRating bool     `json:"decliningRating"`
	Territorialized bool     `json:"territorialized"`
	Active          bool     `json:"active"`
	Favorites       bool     `json:"favorites"`
	UserId          string   `json:"userId"`
	Title           string   `json:"title"`
}

func EventBaseFiltersToDto(in *EventBaseFilters) *dto.EventBaseFilters {
	return &dto.EventBaseFilters{
		Page:            in.Page,
		MaxOnPage:       in.MaxOnPage,
		Tags:            in.Tags,
		DecliningRating: in.DecliningRating,
		Territorialized: in.Territorialized,
		Active:          in.Active,
		Favorites:       in.Favorites,
		UserId:          in.UserId,
		Title:           in.Title,
	}
}

type GetPublicEvent struct {
	EventId         string `json:"id"`
	Title           string `json:"title" validate:"required"`
	Description     string `json:"description"`
	Tags            []Tag  `json:"tags"`
	Cover           string `json:"cover"`
	LastEditionDate string `json:"lastEditionDate"`
	Rate            int64  `json:"rate"`
	Favorite        bool   `json:"favorite"`
}

type GetPublicEventsOut struct {
	Events []*GetPublicEvent `json:"events"`
}

func GetPublicEventFromDto(in *dto.GetPublicEvent) *GetPublicEvent {
	event := &GetPublicEvent{
		EventId:         in.EventId,
		Title:           in.Title,
		Description:     in.Description,
		Tags:            TagsFromDto(in.Tags).Tags,
		Cover:           in.Cover,
		LastEditionDate: in.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000"),
		Rate:            in.Rate,
		Favorite:        in.Favorite,
	}
	if event.LastEditionDate == "01.01.1970 00:00:00.000" {
		event.LastEditionDate = ""
	}

	return event
}

func GetPublicEventsOutFromDto(in *dto.GetPublicEventsOut) *GetPublicEventsOut {
	events := make([]*GetPublicEvent, len(in.Events))
	for i, event := range in.Events {
		events[i] = GetPublicEventFromDto(event)
	}

	return &GetPublicEventsOut{
		Events: events,
	}
}

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
}

func TagsFromDto(in []*dto.Tag) Tags {
	tags := make([]Tag, len(in))
	for i, tag := range in {
		tags[i] = Tag{
			Id:   tag.Id,
			Name: tag.Name,
		}
	}

	return Tags{
		Tags: tags,
	}
}

type EventPlayerInfo struct {
	EventId         string         `json:"eventId"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Tags            []Tag          `json:"tags"`
	Cover           string         `json:"cover"`
	StartDate       string         `json:"startDate"`
	EndDate         string         `json:"endDate"`
	LastEditionDate string         `json:"lastEditionDate"`
	Authors         []Collaborator `json:"authors"`
	Rate            int64          `json:"rate"`
	Favorite        bool           `json:"favorite"`
	Status          string         `json:"status"`
	CanBeDownloaded bool           `json:"canBeDownloaded"`
	IsPrivate       bool           `json:"isPrivate"`
	NeedGroup       bool           `json:"needGroup"`
}

type EventBlockTaskUserIds struct {
	UserId   string `json:"userId"`
	TaskId   string `json:"taskId"`
	BlockId  string `json:"blockId"`
	EventId  string `json:"eventId"`
	Finished bool   `json:"finished"`
}

func EventBlockTaskUserIdsToDto(in *EventBlockTaskUserIds) *dto.EventBlockTaskUserIds {
	return &dto.EventBlockTaskUserIds{
		UserId:   in.UserId,
		TaskId:   in.TaskId,
		BlockId:  in.BlockId,
		EventId:  in.EventId,
		Finished: in.Finished,
	}
}

type UserEventIds struct {
	UserId  string `json:"userId"`
	EventId string `json:"eventId"`
}

func UserEventIdsToDto(in *UserEventIds) *dto.UserEventIds {
	return &dto.UserEventIds{
		UserId:  in.UserId,
		EventId: in.EventId,
	}
}

type PutTimestampIn struct {
	UserId    string `json:"userId"`
	EventId   string `json:"eventId"`
	Timestamp string `json:"timestamp" `
}

func PutTimestampInToDto(in *PutTimestampIn) (*dto.PutTimestampIn, error) {
	ret := &dto.PutTimestampIn{
		UserId:  in.UserId,
		EventId: in.EventId,
	}

	if in.Timestamp != "" {
		ts, err := time.Parse("02.01.2006 15:04:05.000", in.Timestamp)
		if err != nil {
			return nil, err
		}
		ret.Timestamp = timestamppb.New(ts)
	}

	return ret, nil
}

type NextStageInfo struct {
	Type  string          `json:"type"`
	Task  *NextStageTask  `json:"task"`
	Block *NextStageBlock `json:"block"`
}

func NextStageInfoFromDto(in *dto.NextStageInfo) *NextStageInfo {
	return &NextStageInfo{
		Type:  in.Type,
		Task:  NextStageTaskFromDto(in.Task),
		Block: NextStageBlockFromDto(in.Block),
	}
}

type NextStageTask struct {
	TaskId      string       `json:"taskId"`
	BlockId     string       `json:"blockId"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	TaskType    int64        `json:"type"`
	Options     []TaskOption `json:"options"`
	Files       []string     `json:"files"`
	Time        int64        `json:"time"`
	Timestamp   string       `json:"timestamp"`
}

func NextStageTaskFromDto(in *dto.NextStageTask) *NextStageTask {
	if in == nil {
		return nil
	}
	files := make([]string, 0)

	if in.Files != nil || len(in.Files) > 0 {
		files = in.Files
	}

	nextStageTask := &NextStageTask{
		TaskId:      in.TaskId,
		BlockId:     in.BlockId,
		Name:        in.Name,
		Description: in.Description,
		TaskType:    in.Type,
		Options:     TaskOptionsFromDto(in.Options),
		Files:       files,
		Time:        in.Time,
		Timestamp:   in.Timestamp.AsTime().Format("02.01.2006 15:04:05.000"),
	}

	if nextStageTask.Timestamp == "01.01.1970 00:00:00.000" {
		nextStageTask.Timestamp = ""
	}

	return nextStageTask
}

type NextStageBlock struct {
	BlockId    string               `json:"blockId"`
	Name       string               `json:"name"`
	Order      int64                `json:"order"`
	IsParallel bool                 `json:"isParallel"`
	Conditions []Condition          `json:"conditions"`
	Tasks      []NextStageTaskShort `json:"tasks"`
}

func NextStageBlockFromDto(in *dto.NextStageBlock) *NextStageBlock {
	if in == nil {
		return nil
	}
	return &NextStageBlock{
		BlockId:    in.BlockId,
		Name:       in.Name,
		Order:      in.Order,
		IsParallel: in.IsParallel,
		Conditions: ConditionsFromDto(in.Conditions).Conditions,
		Tasks:      NextStageTaskShortsFromDto(in.Tasks),
	}
}

type NextStageTaskShort struct {
	TaskId      string `json:"id"`
	Name        string `json:"name"`
	Type        int64  `json:"type"`
	Time        int64  `json:"time"`
	IsCompleted bool   `json:"isCompleted"`
	// Description string `json:"description"`
}

func NextStageTaskShortsFromDto(in []*dto.NextStageTaskShort) []NextStageTaskShort {
	ret := make([]NextStageTaskShort, len(in))
	for i, task := range in {
		ret[i] = NextStageTaskShort{
			TaskId:      task.TaskId,
			Name:        task.Name,
			Time:        task.Time,
			IsCompleted: task.IsCompleted,
			Type:        task.Type,
			// Description: task.Description,
		}
	}
	return ret
}
