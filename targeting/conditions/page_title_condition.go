package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

type PageTitleCondition struct {
	StringValueCondition
}

func NewPageTitleCondition(c types.TargetingCondition) *PageTitleCondition {
	return &PageTitleCondition{
		StringValueCondition: StringValueCondition{
			TargetingConditionBase: types.TargetingConditionBase{
				Type:    c.Type,
				Include: c.Include,
			},
			Value:     c.Title,
			MatchType: c.MatchType,
		},
	}
}

func (c *PageTitleCondition) CheckTargeting(targetData interface{}) bool {
	targeting := false
	pageViewStorage, ok := targetData.(storage.DataMapStorage[string, types.PageViewVisit])
	if ok && (pageViewStorage != nil) {
		var latest types.PageViewVisit
		pageViewStorage.Enumerate(func(pvv types.PageViewVisit) bool {
			if pvv.LastTimestamp > latest.LastTimestamp {
				latest = pvv
			}
			return true
		})
		return latest.PageView != nil && c.checkTargeting(latest.PageView.Title())
	}
	return targeting
}
