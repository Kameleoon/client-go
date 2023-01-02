package types

type VariationByExposition struct {
	VariationKey string  `json:"variationKey"`
	VariationID  *int    `json:"variationId"`
	Exposition   float64 `json:"exposition"`
}
