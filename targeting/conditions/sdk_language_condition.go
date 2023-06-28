package conditions

import (
	"fmt"

	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
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
	// return false, if sdk language are not equal
	if c.SdkLanguage != sdkInfo.Language {
		return false
	}
	// sdk languages types are equal and version isn't defined - return true
	if len(c.Version) == 0 {
		return true
	}
	// get major / minor / patch sdk version from condition
	majorCondition, minorCondition, patchCondition, err := GetMajorMinorPatch(c.Version)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// get major / minor / patch sdk version from targeted data (provided by SDK)
	majorSdk, minorSdk, patchSdk, err := GetMajorMinorPatch(sdkInfo.Version)
	if err != nil {
		fmt.Println(err)
		return false
	}

	switch c.VersionMatchType {
	case types.OperatorEqual:
		return majorSdk == majorCondition && minorSdk == minorCondition && patchSdk == patchCondition
	case types.OperatorGreater:
		return majorSdk > majorCondition ||
			(majorSdk == majorCondition && minorSdk > minorCondition) ||
			(majorSdk == majorCondition && minorSdk == minorCondition && patchSdk > patchCondition)
	case types.OperatorLower:
		return majorSdk < majorCondition ||
			(majorSdk == majorCondition && minorSdk < minorCondition) ||
			(majorSdk == majorCondition && minorSdk == minorCondition && patchSdk < patchCondition)
	default:
		fmt.Printf("unexpected comparing operation for SdkLanguage condition: %v\n", c.VersionMatchType)
		return false
	}
}

func (c *SdkLanguageCondition) String() string {
	return utils.JsonToString(c)
}
