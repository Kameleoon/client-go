package types

type TrackingTool struct {
	Type                         TrackingToolType `json:"type"`
	CustomVariable               int              `json:"customVariable"`
	GoogleAnalyticsTracker       string           `json:"googleAnalyticsTracker"`
	UniversalAnalyticsDimension  int              `json:"universalAnalyticsDimension"`
	AdobeOmnitureObject          string           `json:"adobeOmnitureObject"`
	EulerianUserCentricParameter string           `json:"eulerianUserCentricParameter"`
	HeatMapPageWidth             int              `json:"heatMapPageWidth"`
	ComScoreCustomerID           string           `json:"comScoreCustomerId"`
	ComScoreDomain               string           `json:"comScoreDomain"`
	ReportingScript              string           `json:"reportingScript"`
}

type TrackingToolType string

const (
	TrackingToolGoogleAnalytics    TrackingToolType = "GOOGLE_ANALYTICS"
	TrackingToolUniversalAnalytics TrackingToolType = "UNIVERSAL_ANALYTICS"
	TrackingToolEconda             TrackingToolType = "ECONDA"
	TrackingToolAtInternet         TrackingToolType = "AT_INTERNET"
	TrackingToolSmartTag           TrackingToolType = "SMART_TAG"
	TrackingToolAdobeOmniture      TrackingToolType = "ADOBE_OMNITURE"
	TrackingToolEulerian           TrackingToolType = "EULERIAN"
	TrackingToolWebtrends          TrackingToolType = "WEBTRENDS"
	TrackingToolHeatmap            TrackingToolType = "HEATMAP"
	TrackingToolKissMetrics        TrackingToolType = "KISS_METRICS"
	TrackingToolPiwik              TrackingToolType = "PIWIK"
	TrackingToolCrazyEgg           TrackingToolType = "CRAZY_EGG"
	TrackingToolComScore           TrackingToolType = "COM_SCORE"
	TrackingToolTealium            TrackingToolType = "TEALIUM"
	TrackingToolYsance             TrackingToolType = "YSANCE"
	TrackingToolMPathy             TrackingToolType = "M_PATHY"
	TrackingToolMandrill           TrackingToolType = "MANDRILL"
	TrackingToolMailperformance    TrackingToolType = "MAILPERFORMANCE"
	TrackingToolSmartfocus         TrackingToolType = "SMARTFOCUS"
	TrackingToolMailjet            TrackingToolType = "MAILJET"
	TrackingToolMailup             TrackingToolType = "MAILUP"
	TrackingToolEmarsys            TrackingToolType = "EMARSYS"
	TrackingToolExpertSender       TrackingToolType = "EXPERT_SENDER"
	TrackingToolTagCommander       TrackingToolType = "TAG_COMMANDER"
	TrackingToolGoogleTagManager   TrackingToolType = "GOOGLE_TAG_MANAGER"
	TrackingToolContentSquare      TrackingToolType = "CONTENT_SQUARE"
	TrackingToolWebtrekk           TrackingToolType = "WEBTREKK"
	TrackingToolCustomIntegration  TrackingToolType = "CUSTOM_INTEGRATION"
	TrackingToolHeap               TrackingToolType = "HEAP"
	TrackingToolSegment            TrackingToolType = "SEGMENT"
	TrackingToolMixpanel           TrackingToolType = "MIXPANEL"
	TrackingToolIabtcf             TrackingToolType = "IABTCF"
	TrackingToolKameleoonTracking  TrackingToolType = "KAMELEOON_TRACKING"
	TrackingToolCustomTracking     TrackingToolType = "CUSTOM_TRACKING"
)
