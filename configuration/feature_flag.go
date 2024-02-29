package configuration

import (
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

func (ff FeatureFlag) GetVariationByKey(key string) (*types.VariationFeatureFlag, bool) {
	for _, v := range ff.Variations {
		if v.Key == key {
			return &v, true
		}
	}
	return nil, false
}

func (ff FeatureFlag) GetVariationKey(varByExp *types.VariationByExposition, rule *Rule) string {
	if varByExp != nil {
		return varByExp.VariationKey
	} else if rule != nil && rule.IsExperimentType() {
		return string(types.VariationOff)
	} else {
		return ff.DefaultVariationKey
	}
}
