package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

type PageViewNumberCondition struct {
	NumberCondition[int]
}

func NewPageViewNumberCondition(c types.TargetingCondition) *PageViewNumberCondition {
	return &PageViewNumberCondition{
		NumberCondition: NumberCondition[int]{
			TargetingConditionBase: types.TargetingConditionBase{
				Type:    c.Type,
				Include: c.Include,
			},
			Value:     c.PageCount,
			MatchType: c.MatchType,
		},
	}
}

func (c *PageViewNumberCondition) CheckTargeting(targetData interface{}) bool {
	pageViewStorage, ok := targetData.(storage.DataMapStorage[string, types.PageViewVisit])
	if ok && (pageViewStorage != nil) {
		var count int
		pageViewStorage.Enumerate(func(pvv types.PageViewVisit) bool {
			count += pvv.Count
			return true
		})
		return c.checkTargeting(count)
	}
	return false
}
