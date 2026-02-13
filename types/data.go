package types

type BaseData interface {
	DataType() DataType
}

type Data interface {
	BaseData

	dataRestriction()
}

type DataType string

const (
	DataTypeAssignedVariation         DataType = "EXPERIMENT"
	DataTypeApplicationVersion        DataType = "APPLICATION_VERSION"
	DataTypeForcedExperimentVariation DataType = "FORCED_EXPERIMENT_VARIATION"
	DataTypeForcedFeatureVariation    DataType = "FORCED_FEATURE_VARIATION"
	DataTypeCustom                    DataType = "CUSTOM"
	DataTypeBrowser                   DataType = "BROWSER"
	DataTypeConversion                DataType = "CONVERSION"
	DataTypeDevice                    DataType = "DEVICE"
	DataTypePageView                  DataType = "PAGE_VIEW"
	DataTypePageViewVisit             DataType = "PAGE_VIEW_VISIT"
	DataTypePersonalization           DataType = "PERSONALIZATION"
	DataTypeUserAgent                 DataType = "USER_AGENT"
	DataTypeCookie                    DataType = "COOKIE"
	DataTypeGeolocation               DataType = "GEOLOCATION"
	DataTypeOperatingSystem           DataType = "OPERATING_SYSTEM"
	DataTypeVisitorVisits             DataType = "VISITOR_VISITS"
	DataTypeKcsHeat                   DataType = "KCS_HEAT"
	DataTypeUniqueIdentifier          DataType = "UNIQUE_IDENTIFIER"
	DataTypeCBScores                  DataType = "CBS"
	DataTypeTargetedSegment           DataType = "TARGETED_SEGMENT"
)
