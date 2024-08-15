package types

import "fmt"

type RespoolTime struct {
	VariationId int     `json:"variationId,string"`
	Value       float64 `json:"value"`
}

func (rt RespoolTime) String() string {
	return fmt.Sprintf(
		"RespoolTime{VariationId:%v,Value:%v}",
		rt.VariationId,
		rt.Value,
	)
}
