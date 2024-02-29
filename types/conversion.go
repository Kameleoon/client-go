package types

import (
	"strconv"

	"github.com/Kameleoon/client-go/v3/utils"
)

const conversionEventType = "conversion"

type Conversion struct {
	duplicationSafeSendableBase
	goalId   int
	revenue  float64
	negative bool
}

func NewConversion(goalId int, negative ...bool) *Conversion {
	return NewConversionWithRevenue(goalId, 0.0, negative...)
}
func NewConversionWithRevenue(goalId int, revenue float64, negative ...bool) *Conversion {
	var negativeValue bool
	if len(negative) > 0 {
		negativeValue = negative[0]
	}
	c := &Conversion{goalId: goalId, revenue: revenue, negative: negativeValue}
	c.initSendale()
	return c
}

func (c *Conversion) dataRestriction() {
	// This method is required to separate external type `Data` from `BaseData` types
}

func (c *Conversion) GoalId() int {
	return c.goalId
}

func (c *Conversion) Revenue() float64 {
	return c.revenue
}

func (c *Conversion) Negative() bool {
	return c.negative
}

func (c *Conversion) QueryEncode() string {
	nonce := c.Nonce()
	if len(nonce) == 0 {
		return ""
	}
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPEventType, conversionEventType)
	qb.Append(utils.QPGoalId, utils.WritePositiveInt(c.goalId))
	qb.Append(utils.QPRevenue, strconv.FormatFloat(c.revenue, 'f', -1, 64))
	qb.Append(utils.QPNegative, strconv.FormatBool(c.negative))
	qb.Append(utils.QPNonce, nonce)
	return qb.String()
}

func (c *Conversion) DataType() DataType {
	return DataTypeConversion
}
