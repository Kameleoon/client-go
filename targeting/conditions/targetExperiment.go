package conditions

import (
	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
)

func NewTargetExperiment(c types.TargetingCondition) *TargetExperiment {
	return &TargetExperiment{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: true,
		},
		ExperimentId: c.Experiment,
		VariationId:  c.Variation,
		Operator:     c.VariationMatchType,
	}
}

type TargetExperiment struct {
	types.TargetingConditionBase
	Operator     types.OperatorType `json:"variationMatchType"`
	ExperimentId int                `json:"experiment"`
	VariationId  int                `json:"variation"`
}

func (c *TargetExperiment) CheckTargeting(targetData interface{}) bool {
	visitorVariationStorage, ok := targetData.(map[int]int)
	if visitorVariationStorage == nil || !ok {
		return false
	}
	targeting := false
	variationStorageForVisitorExist := len(visitorVariationStorage) > 0
	savedVariation, savedVariationExist := visitorVariationStorage[c.ExperimentId]
	switch c.Operator {
	case types.OperatorExact:
		targeting = variationStorageForVisitorExist && savedVariationExist && savedVariation == c.VariationId
	case types.OperatorAny:
		targeting = variationStorageForVisitorExist
	}
	return targeting
}

func (c *TargetExperiment) String() string {
	return utils.JsonToString(c)
}
