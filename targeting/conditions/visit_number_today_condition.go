package conditions

import (
	"time"

	"github.com/Kameleoon/client-go/v3/types"
)

type VisitNumberTodayCondition struct {
	NumberCondition[int]
}

func NewVisitNumberTodayCondition(c types.TargetingCondition) *VisitNumberTodayCondition {
	return &VisitNumberTodayCondition{
		NumberCondition: NumberCondition[int]{
			TargetingConditionBase: types.TargetingConditionBase{
				Type:    c.Type,
				Include: c.Include,
			},
			Value:     c.VisitCount,
			MatchType: c.MatchType,
		},
	}
}

func (c *VisitNumberTodayCondition) CheckTargeting(targetData interface{}) bool {
	td, ok := targetData.(TargetingDataVisitNumberToday)
	if ok && (c.Value != -1) {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
		prevVisits := td.VisitorVisits.PrevVisits()
		var todayVisitNumber int
		for (todayVisitNumber < len(prevVisits)) && (prevVisits[todayVisitNumber].TimeStarted() >= startOfDay) {
			todayVisitNumber++
		}
		if td.CurrentVisitTimeStarted.UnixMilli() >= startOfDay {
			todayVisitNumber++
		}
		return c.checkTargeting(todayVisitNumber)
	}
	return false
}

type TargetingDataVisitNumberToday struct {
	CurrentVisitTimeStarted time.Time
	VisitorVisits           *types.VisitorVisits
}
