package types

type FeatureFlag interface {
	GetId() int
	GetFeatureKey() string
	GetVariations() []VariationFeatureFlag
	GetVariationByKey(key string) (*VariationFeatureFlag, bool)
	GetDefaultVariationKey() string
	GetEnvironmentEnabled() bool
	GetRules() []Rule
	GetMEGroupName() string
}
