package configuration

import "github.com/Kameleoon/client-go/types"

type Configuration struct {
	Settings       types.Settings  `json:"configuration"`
	Experiments    []Experiment    `json:"experiments"`
	FeatureFlags   []FeatureFlag   `json:"featureFlags"`
	FeatureFlagsV2 []FeatureFlagV2 `json:"featureFlagConfigurations"`
}
