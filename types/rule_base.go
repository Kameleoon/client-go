package types

import "github.com/Kameleoon/client-go/v3/errs"

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

func (r *RuleBase) GetVariationByHash(hashDouble float64) *VariationByExposition {
	threshold := hashDouble
	for _, varByExp := range r.VariationByExposition {
		threshold -= varByExp.Exposition
		if threshold < 0 {
			return &varByExp
		}
	}
	return nil
}

func (r *RuleBase) GetVariationByKey(variationKey string) (*VariationByExposition, error) {
	for i := range r.VariationByExposition {
		if r.VariationByExposition[i].VariationKey == variationKey {
			return &r.VariationByExposition[i], nil
		}
	}
	return nil, errs.NewFeatureVariationNotFoundWithVariationKey(r.Id, variationKey)
}

func (r *RuleBase) IsExperimentType() bool {
	return r.Type == RuleTypeExperimentation
}

func (r *RuleBase) IsTargetDeliveryType() bool {
	return r.Type == RuleTypeTargetedDelivery
}
