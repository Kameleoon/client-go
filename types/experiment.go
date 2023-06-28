package types

type ExperimentType string

const (
	ExperimentTypeClassic    ExperimentType = "CLASSIC"
	ExperimentTypeServerSide ExperimentType = "SERVER_SIDE"
	ExperimentTypeDeveloper  ExperimentType = "DEVELOPER"
	ExperimentTypeMVT        ExperimentType = "MVT"
	ExperimentTypeHybrid     ExperimentType = "HYBRID"
)

type Experiment struct {
	ID          int                   `json:"id,string"`
	Name        string                `json:"name"`
	Status      string                `json:"status,omitempty"`
	Segment     Segment               `json:"segment,omitempty"`
	Deviations  []Deviation           `json:"deviations"`
	RespoolTime []RespoolTime         `json:"respoolTime"`
	Variations  []VariationExperiment `json:"variations"`
	SiteEnabled bool                  `json:"siteEnabled,omitempty"`
}
