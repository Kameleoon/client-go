package types

type FeatureFlag interface {
	GetId() int
	GetFeatureKey() string
	GetVariations() []VariationFeatureFlag
	GetDefaultVariationKey() string
	GetEnvironmentEnabled() bool
	GetRules() []Rule
}
