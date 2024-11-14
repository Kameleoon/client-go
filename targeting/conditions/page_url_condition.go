package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

type PageUrlCondition struct {
	StringValueCondition
}

func NewPageUrlCondition(c types.TargetingCondition) *PageUrlCondition {
	return &PageUrlCondition{
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

func (c *PageUrlCondition) CheckTargeting(targetData interface{}) bool {
	pageViewStorage, ok := targetData.(storage.DataMapStorage[string, types.PageViewVisit])
	if ok && (pageViewStorage != nil) {
		var latest types.PageViewVisit
		pageViewStorage.Enumerate(func(pvv types.PageViewVisit) bool {
			if pvv.LastTimestamp > latest.LastTimestamp {
				latest = pvv
			}
			return true
		})
		return latest.PageView != nil && c.checkTargeting(latest.PageView.URL())
	}
	return false
}
