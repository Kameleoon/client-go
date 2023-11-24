package errs

import "fmt"

type FeatureVariationNotFound struct {
	FeatureError
}

func NewFeatureVariationNotFound(featureKey string, variationKey string) *FeatureVariationNotFound {
	msg := fmt.Sprintf("Variation key '%s' not found for feature key '%s'", variationKey, featureKey)
	return &FeatureVariationNotFound{NewFeatureError(msg)}
}
