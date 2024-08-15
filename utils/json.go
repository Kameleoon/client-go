package utils

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"encoding/json"
	"strings"
)

func JsonToString(o interface{}) string {
	if o == nil {
		return ""
	}
	b, err := json.Marshal(o)
	if err != nil {
		logging.Error("Condition can't be parsed to JSON: %s", o)
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
