package targeting

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/targeting/conditions"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

type TargetingManager interface {
	SetDataFile(dataFile types.DataFile)
	CheckTargeting(visitorCode string, campaignId int, expOrFForRule TargetingObjectInterface) bool
}

type targetingManager struct {
	visitorManager storage.VisitorManager
	dataFile       types.DataFile
}

func NewTargetingManager(visitorManager storage.VisitorManager) TargetingManager {
	return &targetingManager{
		visitorManager: visitorManager,
	}
}

func (tm *targetingManager) SetDataFile(dataFile types.DataFile) {
	tm.dataFile = dataFile
}

func (tm *targetingManager) CheckTargeting(
	visitorCode string,
	campaignId int,
	expOrFForRule TargetingObjectInterface,
) bool {
	segment := expOrFForRule.GetTargetingSegment()
	visitor := tm.visitorManager.GetVisitor(visitorCode)
	return segment == nil || segment.CheckTargeting(func(targetingType types.TargetingType) interface{} {
		return tm.getConditionData(targetingType, visitor, visitorCode, campaignId)
	})
}

func (tm *targetingManager) getConditionData(
	targetingType types.TargetingType,
	visitor storage.Visitor,
	visitorCode string,
	campaignId int,
) interface{} {
	var conditionData interface{}
	switch targetingType {
	case types.TargetingCustomDatum:
		if visitor != nil {
			conditionData = visitor.CustomData()
		}
	case types.TargetingBrowser:
		if visitor != nil {
			conditionData = visitor.Browser()
		}
	case types.TargetingDeviceType:
		if visitor != nil {
			conditionData = visitor.Device()
		}
	case types.TargetingPageTitle:
		fallthrough
	case types.TargetingPageUrl:
		fallthrough
	case types.TargetingPageViews:
		fallthrough
	case types.TargetingPreviousPage:
		if visitor != nil {
			conditionData = visitor.PageViewVisits()
		}
	case types.TargetingConversions:
		if visitor != nil {
			conditionData = visitor.Conversions()
		}
	case types.TargetingVisitorCode:
		conditionData = visitorCode
	case types.TargetingSDKLanguage:
		conditionData = &types.TargetedDataSdk{Language: utils.SdkName, Version: utils.SdkVersion}
	case types.TargetingTargetFeatureFlag:
		targetingDataTargetFeatureFlagCondition := conditions.TargetingDataTargetFeatureFlagCondition{
			DataFile: tm.dataFile,
		}
		if visitor != nil {
			targetingDataTargetFeatureFlagCondition.VariationStorage = visitor.Variations()
		}
		conditionData = targetingDataTargetFeatureFlagCondition
	case types.TargetingExclusiveFeatureFlag:
		targetingDataExclusiveFeatureFlag := conditions.TargetingDataExclusiveFeatureFlag{ExperimentId: campaignId}
		if visitor != nil {
			targetingDataExclusiveFeatureFlag.VariationStorage = visitor.Variations()
		}
		conditionData = targetingDataExclusiveFeatureFlag
	case types.TargetingCookie:
		if visitor != nil {
			conditionData = visitor.Cookie()
		}
	case types.TargetingGeolocation:
		if visitor != nil {
			conditionData = visitor.Geolocation()
		}
	case types.TargetingOperatingSystem:
		if visitor != nil {
			conditionData = visitor.OperatingSystem()
		}
	case types.TargetingSegment:
		conditionData = conditions.TargetingDataSegmentCondition{
			DataFile: tm.dataFile,
			TargetingDataGetter: func(targetingType types.TargetingType) interface{} {
				return tm.getConditionData(targetingType, visitor, visitorCode, campaignId)
			},
		}
	case types.TargetingFirstVisit:
		fallthrough
	case types.TargetingLastVisit:
		fallthrough
	case types.TargetingVisits:
		fallthrough
	case types.TargetingSameDayVisits:
		fallthrough
	case types.TargetingNewVisitors:
		if visitor != nil {
			conditionData = visitor.VisitorVisits()
		} else {
			conditionData = (*types.VisitorVisits)(nil)
		}
	case types.TargetingHeatSlice:
		if visitor != nil {
			conditionData = visitor.KcsHeat()
		}
	}
	return conditionData
}
