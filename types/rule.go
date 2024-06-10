package types

type Rule interface {
	GetRuleBase() *RuleBase
	GetSegment() Segment
}
