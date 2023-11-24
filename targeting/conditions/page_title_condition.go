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
		pageViewStorage.Enumerate(func(pvv types.PageViewVisit) bool {
			if c.checkTargeting(pvv.PageView.Title()) {
				targeting = true
				return false
			}
			return true
		})
	}
	return targeting
}
