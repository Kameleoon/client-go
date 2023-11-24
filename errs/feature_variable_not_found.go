package errs

import "fmt"

type FeatureVariableNotFound struct {
	FeatureError
}

func NewFeatureVariableNotFound(featureKey string, variationKey string, variableKey string) *FeatureVariableNotFound {
	msg := fmt.Sprintf("Variable key '%s' not found for variation key '%s' and feature key '%s'",
		variableKey, variationKey, featureKey)
	return &FeatureVariableNotFound{NewFeatureError(msg)}
}
