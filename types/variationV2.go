package types

type VariationDefaults string

const (
	VARIATION_OFF VariationDefaults = "off"
)

type VariationV2 struct {
	Key       string     `json:"key"`
	Variables []Variable `json:"variables"`
}

func (variation VariationV2) GetVariableByKey(key string) (*Variable, bool) {
	for _, variable := range variation.Variables {
		if variable.Key == key {
			return &variable, true
		}
	}
	return nil, false
}
