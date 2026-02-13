package conditions

import (
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewSdkLanguageCondition(c types.TargetingCondition) *SdkLanguageCondition {
	return &SdkLanguageCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		SdkLanguage:      c.SdkLanguage,
		Version:          c.Version,
		VersionMatchType: c.VersionMatchType,
	}
}

type SdkLanguageCondition struct {
	types.TargetingConditionBase
	SdkLanguage      string             `json:"sdkLanguage"`
	Version          string             `json:"version"`
	VersionMatchType types.OperatorType `json:"versionMatchType,omitempty"`
}

func (c *SdkLanguageCondition) CheckTargeting(targetData interface{}) bool {
	sdkInfo, ok := targetData.(*types.TargetedDataSdk)
	return ok && c.checkTargeting(sdkInfo)
}

func (c *SdkLanguageCondition) checkTargeting(sdkInfo *types.TargetedDataSdk) bool {
	return c.SdkLanguage == sdkInfo.Language &&
		(len(c.Version) == 0 || c.versionCondition().CompareWithVersion(sdkInfo.Version))
}

func (c *SdkLanguageCondition) versionCondition() *VersionCondition {
	return &VersionCondition{
		TargetingConditionBase: c.TargetingConditionBase,
		Version:                c.Version,
		VersionMatchType:       c.VersionMatchType,
	}
}

func (c *SdkLanguageCondition) String() string {
	return utils.JsonToString(c)
}
