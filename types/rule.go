package types

import "fmt"

type Rule struct {
	Variations map[string]Variation
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule{Variations:%v}", r.Variations)
}
