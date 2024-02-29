package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
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
	conversionStorage, ok := targetData.(storage.DataCollectionStorage[*types.Conversion])
	if ok && (conversionStorage != nil) {
		if c.GoalId == 0 && conversionStorage.Len() > 0 {
			return true
		}
		targeted := false
		conversionStorage.Enumerate(func(conversion *types.Conversion) bool {
			targeted = c.GoalId == conversion.GoalId()
			return !targeted
		})
		return targeted
	}
	return false
}

func (c *ConversionCondition) String() string {
	return utils.JsonToString(c)
}
