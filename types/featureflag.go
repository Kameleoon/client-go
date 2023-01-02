package types

type FeatureFlag struct {
	ID                int             `json:"id,string"`
	Name              string          `json:"name"`
	IdentificationKey string          `json:"identificationKey"`
	Description       string          `json:"description"`
	Tags              []string        `json:"tags"`
	ExpositionRate    float64         `json:"expositionRate"`
	Segment           Segment         `json:"segment,omitempty"`
	Variations        []Variation     `json:"variations"`
	Goals             []int           `json:"goals"`
	SDKLanguageType   SDKLanguageType `json:"sdkLanguageType"`
	Status            string          `json:"status"`
	DateCreated       TimeNoTZ        `json:"dateCreated"`
	DateModified      TimeNoTZ        `json:"dateModified"`
	RespoolTime       []RespoolTime   `json:"respoolTime"`
	FeatureStatus     string          `json:"featureStatus"`
	Schedules         []Schedule      `json:"schedules"`
	SiteEnabled       bool            `json:"siteEnabled,omitempty"`
}

type Schedule struct {
	DateStart *TimeTZ `json:"dateStart,omitempty"`
	DateEnd   *TimeTZ `json:"dateEnd,omitempty"`
}
