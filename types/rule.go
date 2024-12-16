package types

type Rule interface {
	TargetingObject
	GetVariationByHash(hashDouble float64) *VariationByExposition
	GetVariationByKey(variationKey string) (*VariationByExposition, error)
	IsExperimentType() bool
	IsTargetDeliveryType() bool
	GetRuleBase() *RuleBase
	GetTargetingSegment() Segment
}
