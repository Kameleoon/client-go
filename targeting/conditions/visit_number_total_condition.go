package conditions

import (
	"github.com/Kameleoon/client-go/v3/types"
)

type VisitNumberTotalCondition struct {
	NumberCondition[int]
}

func NewVisitNumberTotalCondition(c types.TargetingCondition) *VisitNumberTotalCondition {
	return &VisitNumberTotalCondition{
		NumberCondition: NumberCondition[int]{
			TargetingConditionBase: types.TargetingConditionBase{
				Type:    c.Type,
				Include: c.Include,
			},
			Value:     c.VisitCount,
			MatchType: c.MatchType,
		},
	}
}

func (c *VisitNumberTotalCondition) CheckTargeting(targetData interface{}) bool {
	vv, ok := targetData.(*types.VisitorVisits)
	if ok && (c.Value != -1) {
		return c.checkTargeting(len(vv.PrevVisits()) + 1) // +1 for current visit
	}
	return false
}
