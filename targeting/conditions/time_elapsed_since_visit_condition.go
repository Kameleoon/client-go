package conditions

import (
	"time"

	"github.com/Kameleoon/client-go/v3/types"
)

type TimeElapsedSinceVisitCondition struct {
	NumberCondition[int64]
}

func NewTimeElapsedSinceVisitCondition(c types.TargetingCondition) *TimeElapsedSinceVisitCondition {
	return &TimeElapsedSinceVisitCondition{
		NumberCondition: NumberCondition[int64]{
			TargetingConditionBase: types.TargetingConditionBase{
				Type:    c.Type,
				Include: c.Include,
			},
			Value:     c.CountInMillis,
			MatchType: c.MatchType,
		},
	}
}

func (c *TimeElapsedSinceVisitCondition) CheckTargeting(targetData interface{}) bool {
	vv, ok := targetData.(*types.VisitorVisits)
	if ok && (vv != nil) && (c.Value != types.UndefinedCountInMillisValue) {
		prevVisits := vv.PrevVisits()
		if len(prevVisits) > 0 {
			now := time.Now().UnixMilli()
			var visitIndex int
			if c.Type == types.TargetingFirstVisit {
				visitIndex = len(prevVisits) - 1
			}
			return c.checkTargeting(now - prevVisits[visitIndex].TimeStarted())
		}
	}
	return false
}
