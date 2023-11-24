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
	targeting := false
	pageViewStorage, ok := targetData.(storage.DataMapStorage[string, types.PageViewVisit])
	if ok && (pageViewStorage != nil) {
		if c.MatchType == types.OperatorExact {
			return pageViewStorage.Get(c.Value).PageView != nil
		}
		pageViewStorage.Enumerate(func(pvv types.PageViewVisit) bool {
			if c.checkTargeting(pvv.PageView.URL()) {
				targeting = true
				return false
			}
			return true
		})
	}
	return targeting
}
