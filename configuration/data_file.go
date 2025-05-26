package configuration

import (
	"fmt"

	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
)

type DataFile struct {
	lastModified                     string
	customDataInfo                   *types.CustomDataInfo
	holdout                          *types.Experiment
	settings                         Settings
	featureFlags                     map[string]*FeatureFlag
	orderedFeatureFlags              []types.FeatureFlag
	meGroups                         map[string]types.MEGroup
	environment                      string
	hasAnyTDRule                     bool
	featureFlagById                  map[int]types.FeatureFlag
	ruleBySegmentId                  map[int]types.Rule
	ruleInfoByExpId                  map[int]types.RuleInfo
	variationById                    map[int]*types.VariationByExposition
	experimentIdsWithJSOrCSSVariable map[int]struct{}
}

func (df DataFile) String() string {
	return fmt.Sprintf(
		"DataFile{environment:'%v',lastModified:'%v',featureFlags:%v,settings:%v}",
		df.environment, df.lastModified, len(df.featureFlags), df.settings,
	)
}

func NewDataFile(configuration Configuration, lastModified string, environment string) *DataFile {
	logging.Debug(
		"CALL: NewDataFile(configuration: %s, lastModified: %s, environment: %s)",
		configuration, lastModified, environment,
	)
	ffs, orderedFFs := collectFeatureFlagsFromConfiguration(configuration)
	featureFlagById, ruleBySegmentId, ruleInfoByExpId, variationById, experimentIdsWithJSOrCSSVariable :=
		collectIndices(ffs)
	cdi := configuration.CustomDataInfo
	if cdi == nil {
		cdi = types.NewCustomDataInfo()
	}
	dataFile := &DataFile{
		lastModified:                     lastModified,
		customDataInfo:                   cdi,
		holdout:                          configuration.Holdout,
		settings:                         configuration.Settings,
		featureFlags:                     ffs,
		orderedFeatureFlags:              orderedFFs,
		meGroups:                         makeMEGroups(orderedFFs),
		environment:                      environment,
		hasAnyTDRule:                     detIfHasAnyTargetedDeliveryRule(ffs),
		featureFlagById:                  featureFlagById,
		ruleBySegmentId:                  ruleBySegmentId,
		ruleInfoByExpId:                  ruleInfoByExpId,
		variationById:                    variationById,
		experimentIdsWithJSOrCSSVariable: experimentIdsWithJSOrCSSVariable,
	}
	logging.Debug(
		"RETURN: NewDataFile(configuration: %s, lastModified: %s, environment: %s)",
		configuration, lastModified, environment,
	)
	return dataFile
}

func collectFeatureFlagsFromConfiguration(
	configuration Configuration,
) (ffs map[string]*FeatureFlag, ordered []types.FeatureFlag) {
	n := len(configuration.FeatureFlags)
	ffs = make(map[string]*FeatureFlag, n)
	ordered = make([]types.FeatureFlag, n)
	for i := 0; i < n; i++ {
		ff := &configuration.FeatureFlags[i]
		ffs[ff.FeatureKey] = ff
		ordered[i] = ff
	}
	return
}

func (df *DataFile) LastModified() string {
	return df.lastModified
}

func (df *DataFile) CustomDataInfo() *types.CustomDataInfo {
	return df.customDataInfo
}

func (df *DataFile) Holdout() *types.Experiment {
	return df.holdout
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

func (df *DataFile) MEGroups() map[string]types.MEGroup {
	return df.meGroups
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

func (df *DataFile) GetRuleInfoByExpId(experimentId int) (types.RuleInfo, bool) {
	ruleInfo, exists := df.ruleInfoByExpId[experimentId]
	return ruleInfo, exists
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
	ruleInfoByExpId map[int]types.RuleInfo,
	variationById map[int]*types.VariationByExposition,
	experimentIdsWithJSOrCSSVariable map[int]struct{},
) {
	featureFlagById = make(map[int]types.FeatureFlag)
	ruleBySegmentId = make(map[int]types.Rule)
	ruleInfoByExpId = make(map[int]types.RuleInfo)
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
			// ruleInfoByExpId
			ruleInfoByExpId[rulePtr.ExperimentId] = types.RuleInfo{FeatureFlag: ff, Rule: rulePtr}
			// variationById
			for iv := len(rulePtr.VariationsByExposition) - 1; iv >= 0; iv-- {
				variationPtr := &rulePtr.VariationsByExposition[iv]
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

func makeMEGroups(featureFlags []types.FeatureFlag) map[string]types.MEGroup {
	meGroupLists := make(map[string][]types.FeatureFlag)
	for _, ff := range featureFlags {
		meGroupName := ff.GetMEGroupName()
		if meGroupName != "" {
			meGroupLists[meGroupName] = append(meGroupLists[meGroupName], ff)
		}
	}
	meGroups := make(map[string]types.MEGroup)
	for meGroupName, meGroupList := range meGroupLists {
		meGroups[meGroupName] = NewMEGroup(meGroupList)
	}
	return meGroups
}
