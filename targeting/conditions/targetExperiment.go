package conditions

import (
	"strings"

	"github.com/Kameleoon/client-go/v2/types"
	"github.com/segmentio/encoding/json"
)

func NewTargetExperiment(c types.TargetingCondition) *TargetExperiment {
	include := false
	if c.Include != nil {
		include = *c.Include
	}
	if c.IsInclude != nil {
		include = *c.IsInclude
	}
	return &TargetExperiment{
		Type:         c.Type,
		Include:      include,
		ExperimentId: c.Experiment,
		VariationId:  c.Variation,
		Operator:     c.VariationMatchType,
	}
}

type TargetExperiment struct {
	Type         types.TargetingType `json:"targetingType"`
	Operator     types.OperatorType  `json:"variationMatchType"`
	Include      bool                `json:"include"`
	ExperimentId int                 `json:"experiment"`
	VariationId  int                 `json:"variation"`
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
	if c == nil {
		return ""
	}
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	var s strings.Builder
	s.Grow(len(b))
	s.Write(b)
	return s.String()
}

func (c TargetExperiment) GetType() types.TargetingType {
	return c.Type
}

func (c *TargetExperiment) SetType(t types.TargetingType) {
	c.Type = t
}

func (c TargetExperiment) GetInclude() bool {
	return c.Include
}

func (c *TargetExperiment) SetInclude(i bool) {
	c.Include = i
}
