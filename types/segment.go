package types

type Segment struct {
	ID                       int             `json:"id"`
	Name                     string          `json:"name"`
	Description              string          `json:"description"`
	ConditionsData           *ConditionsData `json:"conditionsData"`
	SiteID                   int             `json:"siteId"`
	AudienceTracking         bool            `json:"audienceTracking"`
	AudienceTrackingEditable bool            `json:"audienceTrackingEditable"`
	IsFavorite               bool            `json:"isFavorite"`
	DateCreated              TimeNoTZ        `json:"dateCreated"`
	DateModified             TimeNoTZ        `json:"dateModified"`
	Tags                     []string        `json:"tags"`
	ExperimentAmount         int             `json:"experimentAmount,omitempty"`
	PersonalizationAmount    int             `json:"personalizationAmount,omitempty"`
	ExperimentIds            []string        `json:"experiments,omitempty"`
	PersonalizationIds       []string        `json:"personalizations,omitempty"`
}
