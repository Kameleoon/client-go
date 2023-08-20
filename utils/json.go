package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JsonToString(o interface{}) string {
	if o == nil {
		return ""
	}
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Printf("condition can't be parsed to JSON: %v\n", o)
		return ""
	}
	var s strings.Builder
	s.Grow(len(b))
	s.Write(b)
	return s.String()
}

func EscapeJsonStringControlSymbols(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, `\`, `\\`), `"`, `\"`)
}
