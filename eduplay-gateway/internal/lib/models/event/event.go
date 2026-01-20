package eventModel

import (
	dto "eduplay-gateway/internal/generated/clients/event"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostEventIn struct {
	EventId         string   `json:"id"`
	Title           string   `json:"title" validate:"required"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	Cover           string   `json:"cover"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	Private         bool     `json:"private"`
	Password        string   `json:"password"`
	OwnerId         string   `json:"ownerId"`
	LastEditionDate string   `json:"lastEditionDate"`
}

func PostEventInToDto(in *PostEventIn) (*dto.PostEventIn, error) {
	dto := &dto.PostEventIn{
		EventId:     in.EventId,
		Title:       in.Title,
		Description: in.Description,
		Tags:        in.Tags,
		Cover:       in.Cover,
		Private:     in.Private,
		Password:    in.Password,
		OwnerId:     in.OwnerId,
	}

	if in.StartDate != "" {
		startDate, err := time.Parse("02.01.2006 15:04:05.000", in.StartDate)
		if err != nil {
			return nil, err
		}
		dto.StartDate = timestamppb.New(startDate)
	}

	if in.EndDate != "" {
		endDate, err := time.Parse("02.01.2006 15:04:05.000", in.EndDate)
		if err != nil {
			return nil, err
		}
		dto.EndDate = timestamppb.New(endDate)
	}

	if in.LastEditionDate != "" {
		lastEditionDate, err := time.Parse("02.01.2006 15:04:05.000", in.LastEditionDate)
		if err != nil {
			return nil, err
		}
		dto.LastEditionDate = timestamppb.New(lastEditionDate)
	}

	return dto, nil
}

func PostEventInFromDto(in *dto.PostEventIn) *PostEventIn {
	event := &PostEventIn{
		EventId:         in.EventId,
		Title:           in.Title,
		Description:     in.Description,
		Tags:            in.Tags,
		Cover:           in.Cover,
		Private:         in.Private,
		Password:        in.Password,
		OwnerId:         in.OwnerId,
		StartDate:       in.StartDate.AsTime().Format("02.01.2006 15:04:05.000"),
		EndDate:         in.EndDate.AsTime().Format("02.01.2006 15:04:05.000"),
		LastEditionDate: in.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000"),
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

type Group struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Collaborator struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type GetEventSettings struct {
	EventId         string         `json:"id"`
	Title           string         `json:"title" validate:"required"`
	Description     string         `json:"description"`
	Tags            []string       `json:"tags"`
	Cover           string         `json:"cover"`
	StartDate       string         `json:"startDate"`
	EndDate         string         `json:"endDate"`
	Private         bool           `json:"private"`
	Password        string         `json:"password"`
	LastEditionDate string         `json:"lastEditionDate"`
	Groups          []Group        `json:"groups"`
	Rating          bool           `json:"rating"`
	Collaborators   []Collaborator `json:"collaborators"`
	OwnerId         string         `json:"ownerId"`
}

func GetEventSettingsFromDto(event *dto.PostEventIn, groups *dto.GetGroupsOut, collaborators *dto.GetCollaboratorsOut) *GetEventSettings {
	gps := make([]Group, len(groups.Groups))
	for i, group := range groups.Groups {
		gps[i] = Group{
			Login:    group.Login,
			Password: group.Password,
		}
	}

	collabs := make([]Collaborator, len(collaborators.Users))
	for i, user := range collaborators.Users {
		collabs[i] = Collaborator{
			Id:     user.Id,
			Email:  user.Email,
			Avatar: user.Avatar,
		}
	}

	return &GetEventSettings{
		EventId:         event.EventId,
		Title:           event.Title,
		Description:     event.Description,
		Tags:            event.Tags,
		Cover:           event.Cover,
		StartDate:       event.StartDate.AsTime().Format("02.01.2006 15:04:05.000"),
		EndDate:         event.EndDate.AsTime().Format("02.01.2006 15:04:05.000"),
		Private:         event.Private,
		Password:        event.Password,
		LastEditionDate: event.LastEditionDate.AsTime().Format("02.01.2006 15:04:05.000"),
		Groups:          gps,
		Rating:          true,
		Collaborators:   collabs,
		OwnerId:         event.OwnerId,
	}
}

type PostEventBlockIn struct {
	BlockId    string `json:"blockId"`
	EventId    string `json:"eventId"`
	Name       string `json:"name" validate:"required"`
	Order      int    `json:"order" validate:"required"`
	IsParallel bool   `json:"isParallel"`
}

func PostEventBlockToDto(in *PostEventBlockIn) *dto.PostEventBlockIn {
	return &dto.PostEventBlockIn{
		BlockId:    in.BlockId,
		EventId:    in.EventId,
		Name:       in.Name,
		Order:      int64(in.Order),
		IsParallel: in.IsParallel,
	}
}

type BlockCorrectionCheck struct {
	BlockId                string `json:"blockId"`
	Name                   string `json:"name" validate:"required"`
	Order                  int    `json:"order" validate:"required"`
	IsParallel             bool   `json:"isParallel"`
	ConditionWithoutBlocks bool   `json:"conditionWithoutBlocks"`
}

type GetEventBlocksOut struct {
	Name   string
	Blocks []*BlockCorrectionCheck `json:"blocks"`
}

func GetEventBlocksFromDto(in *dto.GetEventBlocksOut) *GetEventBlocksOut {
	eb := make([]*BlockCorrectionCheck, len(in.Blocks))
	for i, block := range in.Blocks {
		condNoBlocks := false
		for _, cond := range block.Conditions {
			if cond.NextBlockId == "" {
				condNoBlocks = true
			}
		}
		eb[i] = &BlockCorrectionCheck{
			BlockId:                block.BlockId,
			Name:                   block.Name,
			Order:                  int(block.Order),
			IsParallel:             block.IsParallel,
			ConditionWithoutBlocks: condNoBlocks,
		}
	}

	return &GetEventBlocksOut{
		Blocks: eb,
	}
}
