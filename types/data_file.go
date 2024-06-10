package types

type DataFile interface {
	GetFeatureFlags() map[string]FeatureFlag
	HasAnyTargetedDeliveryRule() bool
	GetFeatureFlagById(featureFlagId int) FeatureFlag
	GetRuleBySegmentId(segmentId int) Rule
	GetVariation(variationId int) *VariationByExposition
}
