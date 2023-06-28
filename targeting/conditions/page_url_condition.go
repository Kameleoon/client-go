package conditions

import (
	"github.com/Kameleoon/client-go/v2/types"
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
			Value:         c.Url,
			MatchType:     c.MatchType,
			ConditionType: string(types.TargetingPageUrl),
		},
	}
}

func (c *PageUrlCondition) CheckTargeting(targetData interface{}) bool {
	pageView, ok := GetLastTargetingData(targetData, types.DataTypePageView).(*types.PageView)
	return ok && c.checkTargeting(&pageView.URL)
}
