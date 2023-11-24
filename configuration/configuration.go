package configuration

type Configuration struct {
	Settings     Settings      `json:"configuration"`
	FeatureFlags []FeatureFlag `json:"featureFlags"`
}
