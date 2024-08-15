package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewExclusiveFeatureFlagCondition(c types.TargetingCondition) *ExclusiveFeatureFlagCondition {
	return &ExclusiveFeatureFlagCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: true,
		},
	}
}

type ExclusiveFeatureFlagCondition struct {
	types.TargetingConditionBase
}

func (c *ExclusiveFeatureFlagCondition) CheckTargeting(targetData interface{}) bool {
	if targetingData, ok := targetData.(TargetingDataExclusiveFeatureFlag); ok {
		if (targetingData.VariationStorage == nil) || (targetingData.VariationStorage.Len() == 0) {
			return true
		}
		if targetingData.VariationStorage.Len() == 1 {
			return targetingData.VariationStorage.Get(targetingData.ExperimentId) != nil
		}
	}
	return false
}

func (c ExclusiveFeatureFlagCondition) String() string {
	return utils.JsonToString(c)
}

type TargetingDataExclusiveFeatureFlag struct {
	ExperimentId     int
	VariationStorage storage.DataMapStorage[int, *types.AssignedVariation]
}
