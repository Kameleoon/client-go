package errs

import "fmt"

type FeatureExperimentNotFound struct {
	FeatureError
}

func NewFeatureExperimentNotFound(experimentId int) *FeatureExperimentNotFound {
	msg := fmt.Sprintf("Experiment %d is not found", experimentId)
	return &FeatureExperimentNotFound{NewFeatureError(msg)}
}
