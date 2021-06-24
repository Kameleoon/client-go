package types

type URLMatchType string

const (
	URLMatchExact             URLMatchType = "EXACT"
	URLMatchContains          URLMatchType = "CONTAINS"
	URLMatchRegularExpression URLMatchType = "REGULAR_EXPRESSION"
	URLMatchTargetedUrl       URLMatchType = "TARGETED_URL"
)

type AudienceConfigURL struct {
	URL       string       `json:"url"`
	MatchType URLMatchType `json:"matchType"`
}

type SiteType string

const (
	SiteTypeEcommerce SiteType = "ECOMMERCE"
	SiteTypeMedia     SiteType = "MEDIA"
	SiteTypeOther     SiteType = "OTHER"
)

type AudienceConfig struct {
	MainGoal                     int                 `json:"mainGoal"`
	IncludedTargetingTypeList    []TargetingType     `json:"includedTargetingTypeList"`
	ExcludedTargetingTypeList    []TargetingType     `json:"excludedTargetingTypeList"`
	IncludedConfigurationUrlList []AudienceConfigURL `json:"includedConfigurationUrlList"`
	ExcludedConfigurationUrlList []AudienceConfigURL `json:"excludedConfigurationUrlList"`
	IncludedCustomData           []string            `json:"includedCustomData"`
	ExcludedCustomData           []string            `json:"excludedCustomData"`
	IncludedTargetingSegmentList []string            `json:"includedTargetingSegmentList"`
	ExcludedTargetingSegmentList []string            `json:"excludedTargetingSegmentList"`
	SiteType                     SiteType            `json:"siteType"`
	IgnoreURLSettings            bool                `json:"ignoreURLSettings"`
	PredictiveTargeting          bool                `json:"predictiveTargeting"`
	ExcludedGoalList             []string            `json:"excludedGoalList"`
	IncludedExperimentList       []string            `json:"includedExperimentList"`
	ExcludedExperimentList       []string            `json:"excludedExperimentList"`
	IncludedPersonalizationList  []string            `json:"includedPersonalizationList"`
	ExcludedPersonalizationList  []string            `json:"excludedPersonalizationList"`
	CartAmountGoal               int                 `json:"cartAmountGoal"`
	CartAmountValue              int                 `json:"cartAmountValue"`
}
