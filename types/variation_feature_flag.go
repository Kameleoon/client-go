package types

type VariationDefaults string

const (
	VARIATION_OFF VariationDefaults = "off"
)

type VariationFeatureFlag struct {
	Key       string     `json:"key"`
	Variables []Variable `json:"variables"`
}

func (variation VariationFeatureFlag) GetVariableByKey(key string) (*Variable, bool) {
	for _, variable := range variation.Variables {
		if variable.Key == key {
			return &variable, true
		}
	}
	return nil, false
}
