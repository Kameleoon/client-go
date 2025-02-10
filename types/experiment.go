package types

import "github.com/Kameleoon/client-go/v3/errs"

type Experiment struct {
	ExperimentId           int                     `json:"experimentId,omitempty"`
	VariationsByExposition []VariationByExposition `json:"variationByExposition"`
}

func (e *Experiment) GetVariationByHash(hashDouble float64) *VariationByExposition {
	var total float64
	for i := 0; i < len(e.VariationsByExposition); i++ {
		total += e.VariationsByExposition[i].Exposition
		if total >= hashDouble {
			return &e.VariationsByExposition[i]
		}
	}
	return nil
}

func (e *Experiment) GetVariationByKey(variationKey string) (*VariationByExposition, error) {
	for i := range e.VariationsByExposition {
		if e.VariationsByExposition[i].VariationKey == variationKey {
			return &e.VariationsByExposition[i], nil
		}
	}
	return nil, errs.NewFeatureVariationNotFoundWithVariationKeyAndExpId(e.ExperimentId, variationKey)
}
