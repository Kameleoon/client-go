package configuration

import (
	"github.com/Kameleoon/client-go/targeting"
	"github.com/Kameleoon/client-go/types"
	"github.com/segmentio/encoding/json"
)

type Experiment struct {
	types.Experiment
	TargetingSegment *targeting.Segment `json:"-"`
}

func (exp *Experiment) UnmarshalJSON(data []byte) error {
	type ExperimentHidden Experiment
	if err := json.Unmarshal(data, (*ExperimentHidden)(exp)); err != nil {
		return err
	}
	if exp.Segment.ID != 0 {
		exp.TargetingSegment = targeting.NewSegment(&exp.Segment)
	}
	return nil
}

func (exp *Experiment) SiteCodeEnabled() bool {
	return exp.SiteEnabled
}

func (exp *Experiment) GetTargetingSegment() *targeting.Segment {
	return exp.TargetingSegment
}

func (exp *Experiment) GetId() int {
	return exp.ID
}
