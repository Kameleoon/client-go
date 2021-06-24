package types

import (
	"github.com/segmentio/encoding/json"
)

type ExperimentType string

const (
	ExperimentTypeClassic    ExperimentType = "CLASSIC"
	ExperimentTypeServerSide ExperimentType = "SERVER_SIDE"
	ExperimentTypeDeveloper  ExperimentType = "DEVELOPER"
	ExperimentTypeMVT        ExperimentType = "MVT"
	ExperimentTypeHybrid     ExperimentType = "HYBRID"
)

type Experiment struct {
	ID                     int                        `json:"id"`
	SiteID                 int                        `json:"siteId"`
	Name                   string                     `json:"name"`
	BaseURL                string                     `json:"baseURL"`
	Type                   ExperimentType             `json:"type"`
	Description            string                     `json:"description"`
	Tags                   []string                   `json:"tags"`
	TrackingTools          []TrackingTool             `json:"trackingTools"`
	Status                 string                     `json:"status"`
	DateCreated            TimeNoTZ                   `json:"dateCreated"`
	Goals                  []int                      `json:"goals"`
	TargetingSegmentID     int                        `json:"targetingSegmentId"`
	TargetingSegment       interface{}                `json:"-"`
	MainGoalID             int                        `json:"mainGoalId"`
	AutoOptimized          bool                       `json:"autoOptimized"`
	Deviations             Deviations                 `json:"deviations"`
	RespoolTime            RespoolTime                `json:"respoolTime"`
	TargetingConfiguration TargetingConfigurationType `json:"targetingConfiguration"`
	VariationsID           []int                      `json:"variationsId,omitempty"`
	Variations             []Variation                `json:"-"`
	DateModified           TimeNoTZ                   `json:"dateModified"`
	DateStarted            TimeNoTZ                   `json:"dateStarted"`
	DateStatusModified     TimeNoTZ                   `json:"dateStatusModified"`
	IsArchived             bool                       `json:"isArchived"`
	CreatedBy              int                        `json:"createdBy"`
	CommonCssCode          json.RawMessage            `json:"commonCssCode"`
	CommonJavaScriptCode   json.RawMessage            `json:"commonJavaScriptCode"`
}

type Deviations map[string]float64
type RespoolTime map[string]float64

type ExperimentConfig struct {
	IsEditorLaunchedByShortcut     bool              `json:"isEditorLaunchedByShortcut"`
	IsKameleoonReportingEnabled    bool              `json:"isKameleoonReportingEnabled"`
	CustomVariationSelectionScript string            `json:"customVariationSelectionScript"`
	MinWiningReliability           int               `json:"minWiningReliability"`
	AbtestConsent                  ConsentType       `json:"abtestConsent"`
	AbtestConsentOptout            ConsentOptoutType `json:"abtestConsentOptout"`
	BeforeAbtestConsent            BeforeConsentType `json:"beforeAbtestConsent"`
}
