package types

import "fmt"

type Variation struct {
	Key          string
	VariationID  *int
	ExperimentID *int
	Variables    map[string]Variable
}

func (v Variation) IsActive() bool {
	return v.Key == string(VariationOff)
}

func (v Variation) String() string {
	return fmt.Sprintf(
		"Variation{Key:'%v',VariationID:%v,ExperimentID:%v,Variables:%v}",
		v.Key,
		v.VariationID,
		v.ExperimentID,
		v.Variables,
	)
}
