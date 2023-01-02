package configuration

import "github.com/Kameleoon/client-go/types"

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
