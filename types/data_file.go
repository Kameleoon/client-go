package types

type DataFile interface {
	CustomDataInfo() *CustomDataInfo
	Settings() Settings
	GetFeatureFlags() map[string]FeatureFlag
	GetOrderedFeatureFlags() []FeatureFlag
	GetFeatureFlag(featureKey string) (FeatureFlag, error)

	HasAnyTargetedDeliveryRule() bool
	GetFeatureFlagById(featureFlagId int) FeatureFlag
	GetRuleBySegmentId(segmentId int) Rule
	GetRuleInfoByExpId(experimentId int) (RuleInfo, bool)
	GetVariation(variationId int) *VariationByExposition
	HasExperimentJsCssVariable(experimentId int) bool
}

type RuleInfo struct {
	FeatureFlag FeatureFlag
	Rule        Rule
}
