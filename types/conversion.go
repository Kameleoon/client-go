package types

import (
	"strconv"
	"strings"

	"github.com/Kameleoon/client-go/v2/utils"
)

type Conversion struct {
	GoalId   int
	Revenue  float64
	Negative bool
}

func (c Conversion) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=conversion&goalId=")
	b.WriteString(utils.WritePositiveInt(c.GoalId))
	b.WriteString("&revenue=")
	b.WriteString(strconv.FormatFloat(c.Revenue, 'f', -1, 64))
	b.WriteString("&negative=")
	b.WriteString(strconv.FormatBool(c.Negative))
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())
	return b.String()
}

func (c Conversion) DataType() DataType {
	return DataTypeConversion
}
