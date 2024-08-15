package types

type FeatureFlag interface {
	GetId() int
	GetFeatureKey() string
	GetVariations() []VariationFeatureFlag
	GetVariationByKey(key string) (*VariationFeatureFlag, bool)
	GetVariationKey(varByExp *VariationByExposition, rule Rule) string
	GetDefaultVariationKey() string
	GetEnvironmentEnabled() bool
	GetRules() []Rule
}
