package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
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
	targeting := false
	variationStorage, ok := targetData.(storage.DataMapStorage[int, *types.AssignedVariation])
	if ok && (variationStorage != nil) {
		switch c.Operator {
		case types.OperatorExact:
			av := variationStorage.Get(c.ExperimentId)
			targeting = (av != nil) && (av.VariationId() == c.VariationId)
		case types.OperatorAny:
			targeting = variationStorage.Len() > 0
		}
	}
	return targeting
}

func (c *TargetExperiment) String() string {
	return utils.JsonToString(c)
}
