package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

type TargetPersonalizationCondition struct {
	types.TargetingConditionBase
	personalizationId int
}

func NewTargetPersonalizationCondition(c types.TargetingCondition) *TargetPersonalizationCondition {
	return &TargetPersonalizationCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: true,
		},
		personalizationId: c.PersonalizationId,
	}
}

func (c *TargetPersonalizationCondition) CheckTargeting(targetData interface{}) bool {
	if personalizations, ok := targetData.(storage.DataMapStorage[int, *types.Personalization]); ok {
		return personalizations.Get(c.personalizationId) != nil
	}
	return false
}
