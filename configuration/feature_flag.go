package configuration

import (
	"fmt"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
)

type FeatureFlag struct {
	Id                  int                          `json:"id"`
	FeatureKey          string                       `json:"featureKey"`
	Variations          []types.VariationFeatureFlag `json:"variations"`
	DefaultVariationKey string                       `json:"defaultVariationKey"`
	EnvironmentEnabled  bool                         `json:"environmentEnabled"`
	Rules               []Rule                       `json:"rules"`
}

func (ff FeatureFlag) String() string {
	return fmt.Sprintf("FeatureFlag{Id:%v,FeatureKey:'%v',EnvironmentEnabled:%v,DefaultVariationKey:'%v',Rules:%v}",
		ff.Id, ff.FeatureKey, ff.EnvironmentEnabled, ff.DefaultVariationKey, len(ff.Rules))
}

func (ff *FeatureFlag) GetVariationByKey(key string) (*types.VariationFeatureFlag, bool) {
	logging.Debug("CALL: FeatureFlag.GetVariationByKey(key: %s)", key)
	var variation *types.VariationFeatureFlag
	exist := false
	for _, v := range ff.Variations {
		if v.Key == key {
			variation = &v
			exist = true
			break
		}
	}
	logging.Debug("RETURN: FeatureFlag.GetVariationByKey(key: %s) -> (variation: %s, exist: %s)",
		key, variation, exist)
	return variation, exist
}

func (ff *FeatureFlag) GetId() int {
	return ff.Id
}

func (ff *FeatureFlag) GetFeatureKey() string {
	return ff.FeatureKey
}

func (ff *FeatureFlag) GetVariations() []types.VariationFeatureFlag {
	return ff.Variations
}

func (ff *FeatureFlag) GetDefaultVariationKey() string {
	return ff.DefaultVariationKey
}

func (ff *FeatureFlag) GetEnvironmentEnabled() bool {
	return ff.EnvironmentEnabled
}

func (ff *FeatureFlag) GetRules() []types.Rule {
	rules := make([]types.Rule, len(ff.Rules))
	for i := len(ff.Rules) - 1; i >= 0; i-- {
		rules[i] = &ff.Rules[i]
	}
	return rules
}
