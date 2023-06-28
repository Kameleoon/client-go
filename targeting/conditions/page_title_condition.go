package conditions

import (
	"github.com/Kameleoon/client-go/v2/types"
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
			Value:         c.Title,
			MatchType:     c.MatchType,
			ConditionType: string(types.TargetingPageTitle),
		},
	}
}

func (c *PageTitleCondition) CheckTargeting(targetData interface{}) bool {
	pageView, ok := GetLastTargetingData(targetData, types.DataTypePageView).(*types.PageView)
	return ok && c.checkTargeting(&pageView.Title)
}
