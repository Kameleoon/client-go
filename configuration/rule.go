package configuration

import (
	"github.com/Kameleoon/client-go/v2/targeting"
	"github.com/Kameleoon/client-go/v2/types"
	"github.com/segmentio/encoding/json"
)

type Rule struct {
	types.Rule
	TargetingSegment *targeting.Segment `json:"-"`
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	type RuleHidden Rule
	if err := json.Unmarshal(data, (*RuleHidden)(r)); err != nil {
		return err
	}
	if r.Segment.ID != 0 {
		r.TargetingSegment = targeting.NewSegment(&r.Segment)
	}
	return nil
}

func (r *Rule) GetVariationKey(hashDouble float64) *string {
	total := 0.0
	for _, element := range r.VariationByExposition {
		total += element.Exposition
		if total >= hashDouble {
			return &element.VariationKey
		}
	}
	return nil
}

func (r *Rule) GetVariationIdByKey(key string) *int {
	for _, element := range r.VariationByExposition {
		if element.VariationKey == key {
			return element.VariationID
		}
	}
	return nil
}

func (r *Rule) GetTargetingSegment() *targeting.Segment {
	return r.TargetingSegment
}
