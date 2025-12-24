package types

import "fmt"

// Base

type ForcedVariation interface {
	Rule() IRule
	VarByExp() *VariationByExposition
}

type forcedVariation struct {
	rule     IRule
	varByExp *VariationByExposition
}

func (fv *forcedVariation) Rule() IRule {
	return fv.rule
}

func (fv *forcedVariation) VarByExp() *VariationByExposition {
	return fv.varByExp
}

// Feature

type ForcedFeatureVariation struct {
	forcedVariation
	featureKey string
	simulated  bool
}

func NewForcedFeatureVariation(
	featureKey string, rule IRule, varByExp *VariationByExposition, simulated bool,
) *ForcedFeatureVariation {
	return &ForcedFeatureVariation{
		forcedVariation: forcedVariation{rule: rule, varByExp: varByExp},
		featureKey:      featureKey,
		simulated:       simulated,
	}
}

func (ffv *ForcedFeatureVariation) FeatureKey() string {
	return ffv.featureKey
}

func (ffv *ForcedFeatureVariation) Simulated() bool {
	return ffv.simulated
}

func (*ForcedFeatureVariation) DataType() DataType {
	return DataTypeForcedFeatureVariation
}

func (ffv ForcedFeatureVariation) String() string {
	return fmt.Sprintf(
		"ForcedFeatureVariation{featureKey:'%s',rule:%v,varByExp:%v,simulated:%v}",
		ffv.featureKey, ffv.rule, ffv.varByExp, ffv.simulated,
	)
}

// Experiment

type ForcedExperimentVariation struct {
	forcedVariation
	forceTargeting bool
}

func NewForcedExperimentVariation(
	rule IRule, varByExp *VariationByExposition, forceTargeting bool,
) *ForcedExperimentVariation {
	return &ForcedExperimentVariation{
		forcedVariation: forcedVariation{rule: rule, varByExp: varByExp},
		forceTargeting:  forceTargeting,
	}
}

func (fev *ForcedExperimentVariation) ForceTargeting() bool {
	return fev.forceTargeting
}

func (*ForcedExperimentVariation) DataType() DataType {
	return DataTypeForcedExperimentVariation
}

func (fev *ForcedExperimentVariation) String() string {
	return fmt.Sprintf(
		"ForcedExperimentVariation{rule:%v,varByExp:%v,forceTargeting:%v}",
		fev.rule, fev.varByExp, fev.forceTargeting,
	)
}
