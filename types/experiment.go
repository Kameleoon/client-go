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
	SiteID                 int                        `json:"siteId,omitempty"`
	Name                   string                     `json:"name"`
	BaseURL                string                     `json:"baseURL,omitempty"`
	Type                   ExperimentType             `json:"type,omitempty"`
	Description            string                     `json:"description,omitempty"`
	Tags                   []string                   `json:"tags,omitempty"`
	TrackingTools          []TrackingTool             `json:"trackingTools,omitempty"`
	Status                 string                     `json:"status,omitempty"`
	DateCreated            TimeNoTZ                   `json:"dateCreated,omitempty"`
	Goals                  []int                      `json:"goals,omitempty"`
	TargetingSegmentID     int                        `json:"targetingSegmentId,omitempty"`
	TargetingSegment       interface{}                `json:"-"`
	MainGoalID             int                        `json:"mainGoalId,omitempty"`
	AutoOptimized          bool                       `json:"autoOptimized,,omitempty"`
	Deviations             Deviations                 `json:"deviations"`
	RespoolTime            RespoolTime                `json:"respoolTime"`
	TargetingConfiguration TargetingConfigurationType `json:"targetingConfiguration,omitempty"`
	VariationsID           []int                      `json:"variations,omitempty"`
	Variations             []Variation                `json:"-"`
	DateModified           TimeNoTZ                   `json:"dateModified,omitempty"`
	DateStarted            TimeNoTZ                   `json:"dateStarted,omitempty"`
	DateStatusModified     TimeNoTZ                   `json:"dateStatusModified,omitempty"`
	IsArchived             bool                       `json:"isArchived,omitempty"`
	CreatedBy              int                        `json:"createdBy,omitempty"`
	CommonCssCode          json.RawMessage            `json:"commonCssCode,omitempty"`
	CommonJavaScriptCode   json.RawMessage            `json:"commonJavaScriptCode,omitempty"`
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

// GraphQL Helpers
type ExperimentQL struct {
	ID               int             `json:"id,string"`
	Variations       []VariationQL   `json:"variations"`
	Deviations       []DeviationsQL  `json:"deviations"`
	RespoolTime      []RespoolTimeQL `json:"respoolTime"`
	TargetingSegment SegmentQL       `json:"segment"`
	Experiment
}

type RespoolTimeQL struct {
	VariationId string  `json:"variationId"`
	Value       float64 `json:"value"`
}

type DeviationsQL struct {
	VariationId string  `json:"variationId"`
	Value       float64 `json:"value"`
}

type ExperimentDataGraphQL struct {
	Data ExperimentGraphQL `json:"data"`
}
type ExperimentGraphQL struct {
	Experiments EdgeExperimentGraphQL `json:"experiments"`
}
type EdgeExperimentGraphQL struct {
	Edge []NodeExperimentGraphQL `json:"edges"`
}

type NodeExperimentGraphQL struct {
	Node ExperimentQL `json:"node"`
}

func (expNodeQL *NodeExperimentGraphQL) Transform() Experiment {
	expNode := expNodeQL.Node
	exp := expNode.Experiment
	exp.ID = expNode.ID
	//transform variations
	for _, variation := range expNode.Variations {
		exp.VariationsID = append(exp.VariationsID, variation.ID)
		variation.Variation.ID = variation.ID
		exp.Variations = append(exp.Variations, variation.Variation)
	}
	// transform segment
	exp.TargetingSegmentID = expNode.TargetingSegment.ID
	expNode.TargetingSegment.Segment.ID = expNode.TargetingSegment.ID
	exp.TargetingSegment = expNode.TargetingSegment.Segment
	//transform deviations
	exp.Deviations = map[string]float64{}
	for _, deviation := range expNode.Deviations {
		exp.Deviations[deviation.VariationId] = deviation.Value
	}
	//transform respoolTime
	exp.RespoolTime = map[string]float64{}
	for _, respoolTime := range expNode.RespoolTime {
		exp.RespoolTime[respoolTime.VariationId] = respoolTime.Value
	}
	return exp
}
