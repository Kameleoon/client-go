package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

type PreviousPageCondition struct {
	StringValueCondition
}

func NewPreviousPageCondition(c types.TargetingCondition) *PreviousPageCondition {
	return &PreviousPageCondition{
		StringValueCondition: StringValueCondition{
			TargetingConditionBase: types.TargetingConditionBase{
				Type:    c.Type,
				Include: c.Include,
			},
			Value:     c.Url,
			MatchType: c.MatchType,
		},
	}
}

func (c *PreviousPageCondition) CheckTargeting(targetData interface{}) bool {
	pageViewStorage, ok := targetData.(storage.DataMapStorage[string, types.PageViewVisit])
	if ok && (pageViewStorage != nil) {
		var mostRecentVisit, secondMostRecentVisit types.PageViewVisit
		pageViewStorage.Enumerate(func(pvv types.PageViewVisit) bool {
			if (mostRecentVisit.PageView == nil) || (pvv.LastTimestamp > mostRecentVisit.LastTimestamp) {
				secondMostRecentVisit = mostRecentVisit
				mostRecentVisit = pvv
			} else if (secondMostRecentVisit.PageView == nil) ||
				(pvv.LastTimestamp > secondMostRecentVisit.LastTimestamp) {
				secondMostRecentVisit = pvv
			}
			return true
		})
		return (secondMostRecentVisit.PageView != nil) && c.checkTargeting(secondMostRecentVisit.PageView.URL())
	}
	return false
}
