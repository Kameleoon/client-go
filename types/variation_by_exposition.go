package types

import "fmt"

type VariationByExposition struct {
	VariationKey string  `json:"variationKey"`
	VariationID  *int    `json:"variationId"`
	Exposition   float64 `json:"exposition"`
}

func (vbe VariationByExposition) String() string {
	return fmt.Sprintf(
		"VariationByExposition{VariationKey:'%v',VariationID:%v,Exposition:%v}",
		vbe.VariationKey,
		vbe.VariationID,
		vbe.Exposition,
	)
}
