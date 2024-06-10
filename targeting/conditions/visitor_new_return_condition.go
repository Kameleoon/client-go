package conditions

import (
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

const (
	visitorTypeNew    = "NEW"
	visitorTypeReturn = "RETURNING"
)

type VisitorNewReturnCondition struct {
	types.TargetingConditionBase
	visitorType string
}

func NewVisitorNewReturnCondition(c types.TargetingCondition) *VisitorNewReturnCondition {
	return &VisitorNewReturnCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		visitorType: c.VisitorType,
	}
}

func (c *VisitorNewReturnCondition) CheckTargeting(targetData interface{}) bool {
	vv, ok := targetData.(*types.VisitorVisits)
	if ok {
		prevVisitsTime := vv.PreviousVisitTimestamps()
		switch c.visitorType {
		case visitorTypeNew:
			return len(prevVisitsTime) == 0
		case visitorTypeReturn:
			return len(prevVisitsTime) > 0
		}
	}
	return false
}

func (c *VisitorNewReturnCondition) String() string {
	return utils.JsonToString(c)
}
