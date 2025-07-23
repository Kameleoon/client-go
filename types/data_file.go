package types

type DataFile interface {
	LastModified() string
	CustomDataInfo() *CustomDataInfo
	Holdout() *Experiment
	Settings() Settings
	Segments() map[int]Segment
	AudienceTrackingSegments() []Segment
	GetFeatureFlags() map[string]FeatureFlag
	GetOrderedFeatureFlags() []FeatureFlag
	GetFeatureFlag(featureKey string) (FeatureFlag, error)
	MEGroups() map[string]MEGroup

	HasAnyTargetedDeliveryRule() bool
	GetFeatureFlagById(featureFlagId int) FeatureFlag
	GetRuleInfoByExpId(experimentId int) (RuleInfo, bool)
	GetVariation(variationId int) *VariationByExposition
	HasExperimentJsCssVariable(experimentId int) bool
}

type RuleInfo struct {
	FeatureFlag FeatureFlag
	Rule        Rule
}
