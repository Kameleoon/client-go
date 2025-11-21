package errs

import "fmt"

type FeatureEnvironmentDisabled struct {
	FeatureError
}

func NewFeatureEnvironmentDisabled(featureKey string, environment string) *FeatureEnvironmentDisabled {
	if len(environment) == 0 {
		environment = "default environment"
	} else {
		environment = fmt.Sprintf("environment '%s'", environment)
	}
	msg := fmt.Sprintf("Feature '%s' disabled for %s", featureKey, environment)
	return &FeatureEnvironmentDisabled{NewFeatureError(msg)}
}

func NewFeatureEnvironmentDisabledWithMessage(msg string) *FeatureEnvironmentDisabled {
	return &FeatureEnvironmentDisabled{NewFeatureError(msg)}
}
