package conditions

import (
	"strings"

	"github.com/Kameleoon/client-go/v2/types"
	"github.com/segmentio/encoding/json"
)

func NewExclusiveExperiment(c types.TargetingCondition) *ExclusiveExperiment {
	include := true
	return &ExclusiveExperiment{
		Type:    c.Type,
		Include: include,
	}
}

type ExclusiveExperiment struct {
	Type    types.TargetingType `json:"targetingType"`
	Include bool                `json:"include"`
}

func (c *ExclusiveExperiment) CheckTargeting(targetData interface{}) bool {
	if conditionData, ok := targetData.(*types.ExclusiveExperimentTargetedData); ok {
		visitorVariationStorage := conditionData.VisitorVariationStorage
		if len(visitorVariationStorage) == 0 {
			return true
		}
		if _, experimentIdExist := visitorVariationStorage[conditionData.ExperimentId]; experimentIdExist {
			return len(visitorVariationStorage) == 1
		}
	}
	return false
}

func (c *ExclusiveExperiment) String() string {
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

func (c ExclusiveExperiment) GetType() types.TargetingType {
	return c.Type
}

func (c *ExclusiveExperiment) SetType(t types.TargetingType) {
	c.Type = t
}

func (c ExclusiveExperiment) GetInclude() bool {
	return c.Include
}

func (c *ExclusiveExperiment) SetInclude(i bool) {
	c.Include = i
}
