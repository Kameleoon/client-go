package configuration

import (
	"fmt"

	"github.com/Kameleoon/client-go/v3/targeting"
	"github.com/Kameleoon/client-go/v3/types"
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

func (r *Rule) GetTargetingSegment() types.Segment {
	return r.TargetingSegment
}

func (r *Rule) GetRuleBase() *types.RuleBase {
	return &r.RuleBase
}
