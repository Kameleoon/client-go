package errs

import "fmt"

type FeatureVariationNotFound struct {
	FeatureError
}

func NewFeatureVariationNotFound(featureKey string, variationKey string) *FeatureVariationNotFound {
	msg := fmt.Sprintf("Variation key '%s' not found for feature key '%s'", variationKey, featureKey)
	return &FeatureVariationNotFound{NewFeatureError(msg)}
}

func NewFeatureVariationNotFoundWithVariationKey(ruleId int, variationKey string) *FeatureVariationNotFound {
	msg := fmt.Sprintf("Rule %d does not contain variation '%s'", ruleId, variationKey)
	return &FeatureVariationNotFound{NewFeatureError(msg)}
}

func NewFeatureVariationNotFoundWithVariationKeyAndExpId(expId int, variationKey string) *FeatureVariationNotFound {
	msg := fmt.Sprintf("Experiment %d does not contain variation '%s'", expId, variationKey)
	return &FeatureVariationNotFound{NewFeatureError(msg)}
}

func NewFeatureVariationNotFoundWithVariationId(ruleId int, variationId int) *FeatureVariationNotFound {
	msg := fmt.Sprintf("Rule %d does not contain variation %d", ruleId, variationId)
	return &FeatureVariationNotFound{NewFeatureError(msg)}
}
