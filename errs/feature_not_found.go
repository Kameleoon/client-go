package errs

import "fmt"

type FeatureNotFound struct {
	FeatureError
}

func NewFeatureNotFound(featureKey string) *FeatureNotFound {
	msg := fmt.Sprintf("Feature key '%s' not found", featureKey)
	return &FeatureNotFound{NewFeatureError(msg)}
}
