package configuration

import (
	"fmt"

	"github.com/Kameleoon/client-go/v3/types"
)

type Rule struct {
	types.RuleBase
	TargetingSegment types.Segment `json:"-"`
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule{Id:%v}", r.Id)
}

func (r *Rule) applySegments(segments map[int]types.Segment) {
	if r.SegmentId != 0 {
		if segment, exists := segments[r.SegmentId]; exists {
			r.TargetingSegment = segment
		}
	}
}

func (r *Rule) GetTargetingSegment() types.Segment {
	return r.TargetingSegment
}

func (r *Rule) GetRuleBase() *types.RuleBase {
	return &r.RuleBase
}
