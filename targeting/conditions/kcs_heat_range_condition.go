package conditions

import (
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewKcsHeatRangeCondition(c types.TargetingCondition) *KcsHeatRangeCondition {
	return &KcsHeatRangeCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		GoalId:      c.GoalId,
		KeyMomentId: c.KeyMomentId,
		LowerBound:  c.LowerBound,
		UpperBound:  c.UpperBound,
	}
}

type KcsHeatRangeCondition struct {
	types.TargetingConditionBase
	GoalId      int     `json:"goalId"`
	KeyMomentId int     `json:"keyMomentId"`
	LowerBound  float64 `json:"lowerBound"`
	UpperBound  float64 `json:"upperBound"`
}

func (c *KcsHeatRangeCondition) CheckTargeting(targetData interface{}) bool {
	kcsHeat, ok := targetData.(*types.KcsHeat)
	return ok && (kcsHeat != nil) && c.checkTargeting(kcsHeat)
}

func (c *KcsHeatRangeCondition) checkTargeting(kcsHeat *types.KcsHeat) bool {
	if kcsHeat.Values() == nil {
		return false
	}
	if goalScores, ok := kcsHeat.Values()[c.KeyMomentId]; ok {
		var score float64
		score, ok = goalScores[c.GoalId]
		return ok && (score >= c.LowerBound) && (score <= c.UpperBound)
	}
	return false
}

func (c KcsHeatRangeCondition) String() string {
	return utils.JsonToString(c)
}
