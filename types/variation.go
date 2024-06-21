package types

type Variation struct {
	Key          string
	VariationID  *int
	ExperimentID *int
	Variables    map[string]Variable
}
