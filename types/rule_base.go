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
	Experiment
	Order       int         `json:"order"`
	Id          int         `json:"id,omitempty"`
	Type        RuleType    `json:"type"`
	Segment     SegmentBase `json:"segment"`
	Exposition  float64     `json:"exposition"`
	RespoolTime int         `json:"respoolTime,omitempty"`
}

func (r *RuleBase) IsExperimentType() bool {
	return r.Type == RuleTypeExperimentation
}

func (r *RuleBase) IsTargetDeliveryType() bool {
	return r.Type == RuleTypeTargetedDelivery
}
