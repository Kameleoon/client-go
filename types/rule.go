package types

type RuleType string

const (
	RuleTypeExperimentation  RuleType = "EXPERIMENTATION"
	RuleTypeTargetedDelivery RuleType = "TARGETED_DELIVERY"
)

type Rule struct {
	Order                 int                     `json:"order"`
	Type                  string                  `json:"type"`
	Segment               Segment                 `json:"segment"`
	Exposition            float64                 `json:"exposition"`
	ExperimentID          *int                    `json:"experimentId"`
	VariationByExposition []VariationByExposition `json:"variationByExposition"`
}
