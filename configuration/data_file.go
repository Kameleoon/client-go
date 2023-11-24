package configuration

import (
	"github.com/Kameleoon/client-go/v3/errs"
)

type DataFile struct {
	settings     Settings
	featureFlags map[string]FeatureFlag
	environment  string
}

func NewDataFile(configuration Configuration, environment string) *DataFile {
	ffs := make(map[string]FeatureFlag)
	for _, ff := range configuration.FeatureFlags {
		ffs[ff.FeatureKey] = ff
	}
	return &DataFile{
		settings:     configuration.Settings,
		featureFlags: ffs,
	}
}

func (df *DataFile) Settings() Settings {
	return df.settings
}
func (df *DataFile) FeatureFlags() map[string]FeatureFlag {
	return df.featureFlags
}

func (df *DataFile) GetFeatureFlag(featureKey string) (FeatureFlag, error) {
	ff, contains := df.featureFlags[featureKey]
	if !contains {
		return ff, errs.NewFeatureNotFound(featureKey)
	}
	if !ff.EnvironmentEnabled {
		return ff, errs.NewFeatureEnvironmentDisabled(featureKey, df.environment)
	}
	return ff, nil
}

func (df *DataFile) HasAnyTargetedDeliveryRule() bool {
	for _, ff := range df.featureFlags {
		if ff.EnvironmentEnabled {
			for _, rule := range ff.Rules {
				if rule.IsTargetDeliveryType() {
					return true
				}
			}
		}
	}
	return false
}
