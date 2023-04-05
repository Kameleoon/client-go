package types

type RuleType string

const (
	RuleTypeExperimentation  RuleType = "EXPERIMENTATION"
	RuleTypeTargetedDelivery RuleType = "TARGETED_DELIVERY"
)

type Rule struct {
	Order                 int                     `json:"order"`
	Id                    int                     `json:"id,omitempty"`
	Type                  string                  `json:"type"`
	Segment               Segment                 `json:"segment"`
	Exposition            float64                 `json:"exposition"`
	ExperimentId          int                     `json:"experimentId,omitempty"`
	VariationByExposition []VariationByExposition `json:"variationByExposition"`
	RespoolTime           int                     `json:"respoolTime,omitempty"`
}
