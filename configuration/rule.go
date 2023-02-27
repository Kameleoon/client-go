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

func (r *Rule) GetVariationByHash(hashDouble float64) *types.VariationByExposition {
	total := 0.0
	for _, element := range r.VariationByExposition {
		total += element.Exposition
		if total >= hashDouble {
			return &element
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

func (r *Rule) GetVariation(id int) *types.VariationByExposition {
	for _, element := range r.VariationByExposition {
		if *element.VariationID == id {
			return &element
		}
	}
	return nil
}

func (r *Rule) GetTargetingSegment() *targeting.Segment {
	return r.TargetingSegment
}

func (r *Rule) IsExperimentType() bool {
	return r.Type == string(types.RuleTypeExperimentation)
}

func (r *Rule) IsTargetDeliveryType() bool {
	return r.Type == string(types.RuleTypeTargetedDelivery)
}
