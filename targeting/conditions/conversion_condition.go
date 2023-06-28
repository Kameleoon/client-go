package conditions

import (
	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
)

func NewConversionCondition(c types.TargetingCondition) *ConversionCondition {
	return &ConversionCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		GoalId: c.GoalId,
	}
}

type ConversionCondition struct {
	types.TargetingConditionBase
	GoalId int `json:"goalId"`
}

func (c *ConversionCondition) CheckTargeting(targetData interface{}) bool {
	conversion, ok := GetLastTargetingData(targetData, types.DataTypeConversion).(*types.Conversion)
	return ok && (c.GoalId == 0 || c.GoalId == conversion.GoalId)
}

func (c *ConversionCondition) String() string {
	return utils.JsonToString(c)
}
