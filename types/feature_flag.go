package types

import (
	"fmt"
)

type FeatureFlag struct {
	Variations           map[string]Variation
	IsEnvironmentEnabled bool
	Rules                []Rule
	DefaultVariationKey  string
}

func (ff FeatureFlag) DefaultVariation() Variation {
	return ff.Variations[ff.DefaultVariationKey]
}

func (ff FeatureFlag) String() string {
	return fmt.Sprintf(
		"FeatureFlag{Variations:%v,IsEnvironmentEnabled:%v,Rules:%v,DefaultVariationKey:'%v'}",
		ff.Variations, ff.IsEnvironmentEnabled, ff.Rules, ff.DefaultVariationKey,
	)
}
