package types

type RuleType uint8

const (
	RuleTypeUnknown          RuleType = 0
	RuleTypeExperimentation  RuleType = 1
	RuleTypeTargetedDelivery RuleType = 2
)

const (
	ruleTypeLiteralExperimentation  = `"EXPERIMENTATION"`
	ruleTypeLiteralTargetedDelivery = `"TARGETED_DELIVERY"`
)

func (rt *RuleType) UnmarshalJSON(data []byte) error {
	ruleTypeLiteral := string(data)
	switch ruleTypeLiteral {
	case ruleTypeLiteralExperimentation:
		*rt = RuleTypeExperimentation
	case ruleTypeLiteralTargetedDelivery:
		*rt = RuleTypeTargetedDelivery
	default:
		*rt = RuleTypeUnknown
	}
	return nil
}

type RuleBase struct {
	Order                 int                     `json:"order"`
	Id                    int                     `json:"id,omitempty"`
	Type                  RuleType                `json:"type"`
	Segment               SegmentBase             `json:"segment"`
	Exposition            float64                 `json:"exposition"`
	ExperimentId          int                     `json:"experimentId,omitempty"`
	VariationByExposition []VariationByExposition `json:"variationByExposition"`
	RespoolTime           int                     `json:"respoolTime,omitempty"`
}
