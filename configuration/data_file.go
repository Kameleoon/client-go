package configuration

import (
	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/types"
)

type DataFile struct {
	customDataInfo  *types.CustomDataInfo
	settings        Settings
	featureFlags    map[string]*FeatureFlag
	environment     string
	hasAnyTDRule    bool
	featureFlagById map[int]types.FeatureFlag
	ruleBySegmentId map[int]types.Rule
	variationById   map[int]*types.VariationByExposition
}

func NewDataFile(configuration Configuration, environment string) *DataFile {
	ffs := collectFeatureFlagsFromConfiguration(configuration)
	featureFlagById, ruleBySegmentId, variationById := collectIndices(ffs)
	return &DataFile{
		customDataInfo:  configuration.CustomDataInfo,
		settings:        configuration.Settings,
		featureFlags:    ffs,
		environment:     environment,
		hasAnyTDRule:    detIfHasAnyTargetedDeliveryRule(ffs),
		featureFlagById: featureFlagById,
		ruleBySegmentId: ruleBySegmentId,
		variationById:   variationById,
	}
}

func collectFeatureFlagsFromConfiguration(configuration Configuration) map[string]*FeatureFlag {
	ffs := make(map[string]*FeatureFlag)
	for i := len(configuration.FeatureFlags) - 1; i >= 0; i-- {
		ff := &configuration.FeatureFlags[i]
		ffs[ff.FeatureKey] = ff
	}
	return ffs
}

func (df *DataFile) CustomDataInfo() *types.CustomDataInfo {
	return df.customDataInfo
}

func (df *DataFile) Settings() Settings {
	return df.settings
}

func (df *DataFile) FeatureFlags() map[string]*FeatureFlag {
	return df.featureFlags
}

func (df *DataFile) GetFeatureFlag(featureKey string) (*FeatureFlag, error) {
	ff, contains := df.featureFlags[featureKey]
	if !contains {
		return ff, errs.NewFeatureNotFound(featureKey)
	}
	if !ff.EnvironmentEnabled {
		return ff, errs.NewFeatureEnvironmentDisabled(featureKey, df.environment)
	}
	return ff, nil
}

func (df *DataFile) GetFeatureFlags() map[string]types.FeatureFlag {
	ffs := make(map[string]types.FeatureFlag)
	for key, ff := range df.featureFlags {
		ffs[key] = ff
	}
	return ffs
}

func (df *DataFile) HasAnyTargetedDeliveryRule() bool {
	return df.hasAnyTDRule
}

func (df *DataFile) GetFeatureFlagById(featureFlagId int) types.FeatureFlag {
	return df.featureFlagById[featureFlagId]
}

func (df *DataFile) GetRuleBySegmentId(segmentId int) types.Rule {
	return df.ruleBySegmentId[segmentId]
}

func (df *DataFile) GetVariation(variationId int) *types.VariationByExposition {
	return df.variationById[variationId]
}

func detIfHasAnyTargetedDeliveryRule(featureFlags map[string]*FeatureFlag) bool {
	for _, ff := range featureFlags {
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

func collectIndices(featureFlags map[string]*FeatureFlag) (
	featureFlagById map[int]types.FeatureFlag,
	ruleBySegmentId map[int]types.Rule,
	variationById map[int]*types.VariationByExposition,
) {
	featureFlagById = make(map[int]types.FeatureFlag)
	ruleBySegmentId = make(map[int]types.Rule)
	variationById = make(map[int]*types.VariationByExposition)
	for _, ff := range featureFlags {
		for ir := len(ff.Rules) - 1; ir >= 0; ir-- {
			rulePtr := &ff.Rules[ir]
			// ruleBySegmentId
			if rulePtr.TargetingSegment != nil {
				ruleBySegmentId[rulePtr.TargetingSegment.ID] = rulePtr
			}
			// variationById
			for iv := len(rulePtr.VariationByExposition) - 1; iv >= 0; iv-- {
				variationPtr := &rulePtr.VariationByExposition[iv]
				if variationPtr.VariationID != nil {
					variationById[*variationPtr.VariationID] = variationPtr
				}
			}
		}
		// featureFlagById
		featureFlagById[ff.Id] = ff
	}
	return
}
