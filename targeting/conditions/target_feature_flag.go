package conditions

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewTargetFeatureFlagCondition(c types.TargetingCondition) *TargetFeatureFlagCondition {
	return &TargetFeatureFlagCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: true,
		},
		FeatureFlagId:         c.FeatureFlagId,
		ConditionVariationKey: c.VariationKey,
		ConditionRuleId:       c.RuleId,
	}
}

type TargetFeatureFlagCondition struct {
	types.TargetingConditionBase
	FeatureFlagId         int    `json:"featureFlagId,omitempty"`
	ConditionVariationKey string `json:"variationKey,omitempty"`
	ConditionRuleId       int    `json:"ruleId,omitempty"`
}

func (c *TargetFeatureFlagCondition) CheckTargeting(targetData interface{}) bool {
	targetingData, ok := targetData.(TargetingDataTargetFeatureFlagCondition)
	if !ok || (targetingData.DataFile == nil) || (targetingData.VariationStorage == nil) ||
		(targetingData.VariationStorage.Len() == 0) {
		return false
	}
	for _, rule := range c.getRules(targetingData.DataFile) {
		if (rule == nil) || (rule.GetRuleBase() == nil) {
			continue
		}
		assignedVariation := targetingData.VariationStorage.Get(rule.GetRuleBase().ExperimentId)
		if assignedVariation == nil {
			continue
		}
		if c.ConditionVariationKey == "" {
			return true
		}
		variation := targetingData.DataFile.GetVariation(assignedVariation.VariationId())
		if (variation != nil) && (variation.VariationKey == c.ConditionVariationKey) {
			return true
		}
	}
	return false
}

func (c *TargetFeatureFlagCondition) getRules(dataFile types.DataFile) []types.Rule {
	ff := dataFile.GetFeatureFlagById(c.FeatureFlagId)
	if ff == nil {
		return nil
	}
	rules := ff.GetRules()
	if c.ConditionRuleId > 0 {
		for _, rule := range rules {
			if (rule.GetRuleBase() != nil) && (rule.GetRuleBase().Id == c.ConditionRuleId) {
				return []types.Rule{rule}
			}
		}
		return nil
	} else {
		return rules
	}
}

func (c TargetFeatureFlagCondition) String() string {
	return utils.JsonToString(c)
}

type TargetingDataTargetFeatureFlagCondition struct {
	DataFile         types.DataFile
	VariationStorage storage.DataMapStorage[int, *types.AssignedVariation]
}
