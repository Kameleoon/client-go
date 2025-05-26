package types

type DataFile interface {
	LastModified() string
	CustomDataInfo() *CustomDataInfo
	Holdout() *Experiment
	Settings() Settings
	GetFeatureFlags() map[string]FeatureFlag
	GetOrderedFeatureFlags() []FeatureFlag
	GetFeatureFlag(featureKey string) (FeatureFlag, error)
	MEGroups() map[string]MEGroup

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
