package types

import "fmt"

// Base

type ForcedVariation struct {
	rule     Rule
	varByExp *VariationByExposition
}

func (fv *ForcedVariation) Rule() Rule {
	return fv.rule
}

func (fv *ForcedVariation) VarByExp() *VariationByExposition {
	return fv.varByExp
}

// Feature

type ForcedFeatureVariation struct {
	ForcedVariation
	featureKey string
	simulated  bool
}

func NewForcedFeatureVariation(
	featureKey string, rule Rule, varByExp *VariationByExposition, simulated bool,
) *ForcedFeatureVariation {
	return &ForcedFeatureVariation{
		ForcedVariation: ForcedVariation{rule: rule, varByExp: varByExp},
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
	ForcedVariation
	forceTargeting bool
}

func NewForcedExperimentVariation(
	rule Rule, varByExp *VariationByExposition, forceTargeting bool,
) *ForcedExperimentVariation {
	return &ForcedExperimentVariation{
		ForcedVariation: ForcedVariation{rule: rule, varByExp: varByExp},
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
