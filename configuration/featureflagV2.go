package configuration

import (
	"github.com/Kameleoon/client-go/v2/types"
)

type FeatureFlagV2 struct {
	ID                  int                 `json:"id"`
	FeatureKey          string              `json:"featureKey"`
	Variations          []types.VariationV2 `json:"variations"`
	DefaultVariationKey string              `json:"defaultVariationKey"`
	Rules               []Rule              `json:"rules"`
}

func (ff FeatureFlagV2) GetVariationByKey(key string) (*types.VariationV2, bool) {
	for _, v := range ff.Variations {
		if v.Key == key {
			return &v, true
		}
	}
	return nil, false
}

func (ff FeatureFlagV2) GetVariationKey(varByExp *types.VariationByExposition, rule *Rule) string {
	if varByExp != nil {
		return varByExp.VariationKey
	} else if rule != nil && rule.IsExperimentType() {
		return string(types.VARIATION_OFF)
	} else {
		return ff.DefaultVariationKey
	}
}
