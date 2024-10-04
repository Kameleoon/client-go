package configuration

import (
	"fmt"
	"sort"

	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
)

type DataFile struct {
	customDataInfo                   *types.CustomDataInfo
	settings                         Settings
	featureFlags                     map[string]*FeatureFlag
	orderedFeatureFlags              []types.FeatureFlag
	environment                      string
	hasAnyTDRule                     bool
	featureFlagById                  map[int]types.FeatureFlag
	ruleBySegmentId                  map[int]types.Rule
	variationById                    map[int]*types.VariationByExposition
	experimentIdsWithJSOrCSSVariable map[int]struct{}
}

func (df DataFile) String() string {
	return fmt.Sprintf("DataFile{environment:'%v',featureFlags:%v,settings:%v}",
		df.environment, len(df.featureFlags), df.settings)
}

func NewDataFile(configuration Configuration, environment string) *DataFile {
	logging.Debug("CALL: NewDataFile(configuration: %s, environment: %s)",
		configuration, environment)
	ffs := collectFeatureFlagsFromConfiguration(configuration)
	orderedFFs := collectOrderedFeatureFlags(ffs)
	featureFlagById, ruleBySegmentId, variationById, experimentIdsWithJSOrCSSVariable := collectIndices(ffs)
	cdi := configuration.CustomDataInfo
	if cdi == nil {
		cdi = types.NewCustomDataInfo()
	}
	dataFile := &DataFile{
		customDataInfo:                   cdi,
		settings:                         configuration.Settings,
		featureFlags:                     ffs,
		orderedFeatureFlags:              orderedFFs,
		environment:                      environment,
		hasAnyTDRule:                     detIfHasAnyTargetedDeliveryRule(ffs),
		featureFlagById:                  featureFlagById,
		ruleBySegmentId:                  ruleBySegmentId,
		variationById:                    variationById,
		experimentIdsWithJSOrCSSVariable: experimentIdsWithJSOrCSSVariable,
	}
	logging.Debug("RETURN: NewDataFile(configuration: %s, environment: %s) -> (dataFile: %s)",
		configuration, environment, dataFile)
	return dataFile
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

func (df *DataFile) Settings() types.Settings {
	return df.settings
}

func (df *DataFile) FeatureFlags() map[string]*FeatureFlag {
	return df.featureFlags
}

func (df *DataFile) GetFeatureFlag(featureKey string) (types.FeatureFlag, error) {
	logging.Debug("CALL: DataFile.GetFeatureFlag(featureKey: %s)", featureKey)
	ff, contains := df.featureFlags[featureKey]
	var err error
	if !contains {
		err = errs.NewFeatureNotFound(featureKey)
	} else if !ff.EnvironmentEnabled {
		err = errs.NewFeatureEnvironmentDisabled(featureKey, df.environment)
	}
	logging.Debug("RETURN: DataFile.GetFeatureFlag(featureKey: %s) -> (featureFlag: %s, error: %s)",
		featureKey, ff, err)
	return ff, err
}

func (df *DataFile) GetFeatureFlags() map[string]types.FeatureFlag {
	logging.Debug("CALL: DataFile.GetFeatureFlags()")
	ffs := make(map[string]types.FeatureFlag)
	for key, ff := range df.featureFlags {
		ffs[key] = ff
	}
	logging.Debug("RETURN: DataFile.GetFeatureFlags() -> (featureFlags: %s)", ffs)
	return ffs
}

func (df *DataFile) GetOrderedFeatureFlags() []types.FeatureFlag {
	return df.orderedFeatureFlags
}

func collectOrderedFeatureFlags(ffs map[string]*FeatureFlag) []types.FeatureFlag {
	ordered := make([]types.FeatureFlag, 0, len(ffs))
	for _, ff := range ffs {
		ordered = append(ordered, ff)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].GetId() < ordered[j].GetId()
	})
	return ordered
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

func (df *DataFile) HasExperimentJsCssVariable(experimentId int) bool {
	_, exists := df.experimentIdsWithJSOrCSSVariable[experimentId]
	return exists
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
	experimentIdsWithJSOrCSSVariable map[int]struct{},
) {
	featureFlagById = make(map[int]types.FeatureFlag)
	ruleBySegmentId = make(map[int]types.Rule)
	variationById = make(map[int]*types.VariationByExposition)
	experimentIdsWithJSOrCSSVariable = make(map[int]struct{})
	for _, ff := range featureFlags {
		hasFeatureFlagVariableJsCss := hasFeatureFlagVariableJsCss(ff)
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
			// experimentIdsWithJSOrCSSVariable
			if hasFeatureFlagVariableJsCss {
				experimentIdsWithJSOrCSSVariable[rulePtr.ExperimentId] = struct{}{}
			}
		}
		// featureFlagById
		featureFlagById[ff.Id] = ff
	}
	return
}

func hasFeatureFlagVariableJsCss(featureFlag *FeatureFlag) bool {
	if len(featureFlag.GetVariations()) > 0 {
		variation := featureFlag.GetVariations()[0]
		for _, variable := range variation.Variables {
			if variable.Type == "JS" || variable.Type == "CSS" {
				return true
			}
		}
	}
	return false
}
