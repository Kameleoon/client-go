package types

type Rule interface {
	TargetingObject
	GetVariationByHash(hashDouble float64) *VariationByExposition
	IsExperimentType() bool
	IsTargetDeliveryType() bool
	GetRuleBase() *RuleBase
	GetTargetingSegment() Segment
}
