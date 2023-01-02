package types

import (
	"strconv"

	"github.com/segmentio/encoding/json"
)

type Deviation struct {
	VariationId int     `json:"variationId,string"`
	Value       float64 `json:"value"`
}

func (deviation *Deviation) UnmarshalJSON(data []byte) error {

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if variationIdInt, err := strconv.Atoi(v["variationId"].(string)); err == nil {
		deviation.VariationId = variationIdInt
	} else {
		deviation.VariationId = 0
	}
	deviation.Value, _ = v["value"].(float64)
	return nil
}
