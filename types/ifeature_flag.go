package types

type IFeatureFlag interface {
	GetId() int
	GetFeatureKey() string
	GetVariations() []VariationFeatureFlag
	GetVariationByKey(key string) (*VariationFeatureFlag, bool)
	GetDefaultVariationKey() string
	GetEnvironmentEnabled() bool
	GetRules() []IRule
	GetMEGroupName() string
	GetBucketingCustomDataIndex() *int
}
