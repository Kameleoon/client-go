package configuration

import (
	"github.com/Kameleoon/client-go/v3/targeting"
	"github.com/Kameleoon/client-go/v3/types"
	"fmt"
	"github.com/segmentio/encoding/json"
)

type Rule struct {
	types.RuleBase
	TargetingSegment *targeting.Segment `json:"-"`
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule{Id:%v}", r.Id)
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
	threshold := hashDouble
	for _, varByExp := range r.VariationByExposition {
		threshold -= varByExp.Exposition
		if threshold < 0 {
			return &varByExp
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

func (r *Rule) GetTargetingSegment() types.Segment {
	return r.TargetingSegment
}

func (r *Rule) IsExperimentType() bool {
	return r.Type == types.RuleTypeExperimentation
}

func (r *Rule) IsTargetDeliveryType() bool {
	return r.Type == types.RuleTypeTargetedDelivery
}

func (r *Rule) GetRuleBase() *types.RuleBase {
	return &r.RuleBase
}
