package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
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
	if conditionData, ok := targetData.(TargetedDataExclusiveExperiment); ok {
		if conditionData.VariationStorage.Len() == 0 {
			return true
		}
		if conditionData.VariationStorage.Len() == 1 {
			return conditionData.VariationStorage.Get(conditionData.ExperimentId) != nil
		}
	}
	return false
}

func (c *ExclusiveExperiment) String() string {
	return utils.JsonToString(c)
}

type TargetedDataExclusiveExperiment struct {
	ExperimentId     int
	VariationStorage storage.DataMapStorage[int, *types.AssignedVariation]
}
