package types

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"fmt"
)

type Variable struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (v Variable) String() string {
	return fmt.Sprintf(
		"Variable{Key:'%v',Type:'%v',Value:%v}",
		v.Key,
		v.Type,
		logging.ObjectToString(v.Value),
	)
}
