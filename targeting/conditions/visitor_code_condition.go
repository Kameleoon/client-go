package conditions

import (
	"github.com/Kameleoon/client-go/v3/types"
)

func NewVisitorCodeCondition(c types.TargetingCondition) *StringValueCondition {
	return &StringValueCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		Value:         c.VisitorCode,
		MatchType:     c.MatchType,
		ConditionType: string(types.TargetingVisitorCode),
	}
}
