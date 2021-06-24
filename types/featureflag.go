package types

type FeatureFlag struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	IdentificationKey  string          `json:"identificationKey"`
	Description        string          `json:"description"`
	Tags               []string        `json:"tags"`
	SiteID             int             `json:"siteId"`
	ExpositionRate     float64         `json:"expositionRate"`
	TargetingSegmentID int             `json:"targetingSegmentId"`
	TargetingSegment   interface{}     `json:"targetingSegment,omitempty"`
	VariationsID       []int           `json:"variationsId,omitempty"`
	Variations         []Variation     `json:"variations,omitempty"`
	Goals              []int           `json:"goals"`
	SDKLanguageType    SDKLanguageType `json:"sdkLanguageType"`
	Status             string          `json:"status"`
	DateCreated        TimeNoTZ        `json:"dateCreated"`
	DateModified       TimeNoTZ        `json:"dateModified"`
}
