package types

type IDataFile interface {
	LastModified() string
	CustomDataInfo() *CustomDataInfo
	Holdout() *Experiment
	Settings() Settings
	Segments() map[int]Segment
	AudienceTrackingSegments() []Segment
	GetFeatureFlags() map[string]IFeatureFlag
	GetOrderedFeatureFlags() []IFeatureFlag
	GetFeatureFlag(featureKey string) (IFeatureFlag, error)
	MEGroups() map[string]MEGroup

	HasAnyTargetedDeliveryRule() bool
	GetFeatureFlagById(featureFlagId int) IFeatureFlag
	GetRuleInfoByExpId(experimentId int) (RuleInfo, bool)
	GetVariation(variationId int) *VariationByExposition
	HasExperimentJsCssVariable(experimentId int) bool
}

type RuleInfo struct {
	FeatureFlag IFeatureFlag
	Rule        IRule
}
