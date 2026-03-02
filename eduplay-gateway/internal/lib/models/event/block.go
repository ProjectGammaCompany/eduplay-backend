package eventModel

import (
	dto "eduplay-gateway/internal/generated/clients/event"
)

type EventBlockName struct {
	BlockId string `json:"Id"`
	Name    string `json:"name" validate:"required"`
}

type PostEventBlockIn struct {
	BlockId       string `json:"id"`
	EventId       string `json:"eventId"`
	Name          string `json:"name" validate:"required"`
	Order         int    `json:"order" validate:"required"`
	IsParallel    bool   `json:"isParallel"`
	Points        bool   `json:"points"`
	Answers       bool   `json:"answers"`
	PartialPoints bool   `json:"partialPoints"`
}

type PutEventBlockIn struct {
	BlockId       string `json:"id"`
	EventId       string `json:"eventId"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	IsParallel    bool   `json:"isParallel"`
	Points        bool   `json:"points"`
	Answers       bool   `json:"answers"`
	PartialPoints bool   `json:"partialPoints"`
}

func PutEventBlockToDto(in *PutEventBlockIn) *dto.PostEventBlockIn {
	return &dto.PostEventBlockIn{
		BlockId:       in.BlockId,
		EventId:       in.EventId,
		Name:          in.Name,
		Order:         int64(in.Order),
		IsParallel:    in.IsParallel,
		ShowPoints:    in.Points,
		ShowAnswers:   in.Answers,
		PartialPoints: in.PartialPoints,
	}
}

func PostEventBlockToDto(in *PostEventBlockIn) *dto.PostEventBlockIn {
	return &dto.PostEventBlockIn{
		BlockId:       in.BlockId,
		EventId:       in.EventId,
		Name:          in.Name,
		Order:         int64(in.Order),
		IsParallel:    in.IsParallel,
		ShowPoints:    in.Points,
		ShowAnswers:   in.Answers,
		PartialPoints: in.PartialPoints,
	}
}

func PostEventBlockFromDto(in *dto.PostEventBlockIn) *PostEventBlockIn {
	return &PostEventBlockIn{
		BlockId:       in.BlockId,
		EventId:       in.EventId,
		Name:          in.Name,
		Order:         int(in.Order),
		IsParallel:    in.IsParallel,
		Points:        in.ShowPoints,
		Answers:       in.ShowAnswers,
		PartialPoints: in.PartialPoints,
	}
}

type BlockCorrectionCheck struct {
	BlockId                string `json:"id"`
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

type Condition struct {
	СonditionId     string   `json:"id"`
	PreviousBlockId string   `json:"previousBlockId"`
	NextBlockId     string   `json:"blockId"`
	NextBlockOrder  int64    `json:"blockOrder"`
	GroupIds        []string `json:"group"`
	Min             int64    `json:"min"`
	Max             int64    `json:"max"`
}

func ConditionToDto(in *Condition) *dto.Condition {
	return &dto.Condition{
		ConditionId:     in.СonditionId,
		PreviousBlockId: in.PreviousBlockId,
		NextBlockId:     in.NextBlockId,
		NextBlockOrder:  in.NextBlockOrder,
		GroupIds:        in.GroupIds,
		Min:             in.Min,
		Max:             in.Max,
	}
}

func ConditionFromDto(in *dto.Condition) *Condition {
	if len(in.GroupIds) == 0 {
		in.GroupIds = []string{}
	}

	return &Condition{
		СonditionId:     in.ConditionId,
		PreviousBlockId: in.PreviousBlockId,
		NextBlockId:     in.NextBlockId,
		NextBlockOrder:  in.NextBlockOrder,
		GroupIds:        in.GroupIds,
		Min:             in.Min,
		Max:             in.Max,
	}
}

type Conditions struct {
	Conditions []Condition `json:"conditions"`
}

func ConditionsFromDto(in []*dto.Condition) *Conditions {
	conditions := make([]Condition, 0)

	for _, condition := range in {
		conditions = append(conditions, *ConditionFromDto(condition))
	}

	return &Conditions{Conditions: conditions}
}

type PostConditionOut struct {
	BlockOrder  int64  `json:"blockOrder"`
	ConditionId string `json:"conditionId"`
}

type GetBlock struct {
	BlockId    string `json:"id"`
	Name       string `json:"name"`
	Order      int    `json:"order"`
	IsParallel bool   `json:"isParallel"`
}

type GetBlockForConditionsOut struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type GetBlocksForConditionsOut struct {
	Blocks []*GetBlockForConditionsOut `json:"blocks"`
}

func GetBlockForConditionsOutFromDto(in *dto.BlockInfo) *GetBlockForConditionsOut {
	return &GetBlockForConditionsOut{
		Id:   in.BlockId,
		Name: in.Name,
	}
}

func GetBlocksForConditionsOutFromDto(in *dto.GetEventBlocksOut) *GetBlocksForConditionsOut {
	eb := make([]*GetBlockForConditionsOut, len(in.Blocks))
	for i, block := range in.Blocks {
		eb[i] = GetBlockForConditionsOutFromDto(block)
	}

	return &GetBlocksForConditionsOut{
		Blocks: eb,
	}
}
