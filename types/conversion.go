package types

import (
	"strconv"

	"github.com/Kameleoon/client-go/v2/network"
	"github.com/Kameleoon/client-go/v2/utils"
)

const conversionEventType = "conversion"

type Conversion struct {
	GoalId   int
	Revenue  float64
	Negative bool
}

func (c Conversion) QueryEncode() string {
	qb := network.NewQueryBuilder()
	qb.Append(network.QPEventType, conversionEventType)
	qb.Append(network.QPGoalId, utils.WritePositiveInt(c.GoalId))
	qb.Append(network.QPRevenue, strconv.FormatFloat(c.Revenue, 'f', -1, 64))
	qb.Append(network.QPNegative, strconv.FormatBool(c.Negative))
	qb.Append(network.QPNonce, network.GetNonce())
	return qb.String()
}

func (c Conversion) DataType() DataType {
	return DataTypeConversion
}
