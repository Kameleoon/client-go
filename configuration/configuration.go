package configuration

import "github.com/Kameleoon/client-go/v2/types"

type Configuration struct {
	Settings       types.Settings  `json:"configuration"`
	Experiments    []Experiment    `json:"experiments"`
	FeatureFlags   []FeatureFlag   `json:"featureFlags"`
	FeatureFlagsV2 []FeatureFlagV2 `json:"featureFlagConfigurations"`
}
