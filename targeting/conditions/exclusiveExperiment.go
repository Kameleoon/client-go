package conditions

import (
	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
)

func NewExclusiveExperiment(c types.TargetingCondition) *ExclusiveExperiment {
	return &ExclusiveExperiment{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: true,
		},
	}
}

type ExclusiveExperiment struct {
	types.TargetingConditionBase
}

func (c *ExclusiveExperiment) CheckTargeting(targetData interface{}) bool {
	if conditionData, ok := targetData.(*types.TargetedDataExclusiveExperiment); ok {
		visitorVariationStorage := conditionData.VisitorVariationStorage
		if len(visitorVariationStorage) == 0 {
			return true
		}
		if len(visitorVariationStorage) == 1 {
			_, experimentIdFound := visitorVariationStorage[conditionData.ExperimentId]
			return experimentIdFound
		}
	}
	return false
}

func (c *ExclusiveExperiment) String() string {
	return utils.JsonToString(c)
}
