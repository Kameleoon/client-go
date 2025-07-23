package conditions

import (
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewSegmentCondition(c types.TargetingCondition) *SegmentCondition {
	return &SegmentCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		SegmentId: c.SegmentId,
	}
}

type SegmentCondition struct {
	types.TargetingConditionBase
	SegmentId int `json:"segmentId,omitempty"`
}

func (c *SegmentCondition) CheckTargeting(targetData interface{}) bool {
	targetingData, ok := targetData.(TargetingDataSegmentCondition)
	if !ok || (targetingData.DataFile == nil) || (targetingData.TargetingDataGetter == nil) {
		return false
	}
	segment := targetingData.DataFile.Segments()[c.SegmentId]
	if segment == nil {
		return false
	}
	return segment.CheckTargeting(targetingData.TargetingDataGetter)
}

func (c SegmentCondition) String() string {
	return utils.JsonToString(c)
}

type TargetingDataSegmentCondition struct {
	DataFile            types.DataFile
	TargetingDataGetter types.TargetingDataGetter
}
