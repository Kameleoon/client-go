package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Kameleoon/client-go/v3/utils"
)

const conversionEventType = "conversion"

type Conversion struct {
	duplicationSafeSendableBase
	goalId   int
	revenue  float64
	negative bool
	metadata []*CustomData
}

type ConversionOptParams struct {
	Revenue  float64
	Negative bool
	Metadata []*CustomData
}

func NewConversion(goalId int, negative ...bool) *Conversion {
	var p ConversionOptParams
	if len(negative) > 0 {
		p.Negative = negative[0]
	}
	return NewConversionWithOptParams(goalId, p)
}
func NewConversionWithRevenue(goalId int, revenue float64, negative ...bool) *Conversion {
	p := ConversionOptParams{Revenue: revenue}
	if len(negative) > 0 {
		p.Negative = negative[0]
	}
	return NewConversionWithOptParams(goalId, p)
}
func NewConversionWithOptParams(goalId int, params ConversionOptParams) *Conversion {
	c := &Conversion{
		goalId:   goalId,
		revenue:  params.Revenue,
		negative: params.Negative,
		metadata: params.Metadata,
	}
	c.initSendale()
	return c
}

func (c Conversion) String() string {
	return fmt.Sprintf("Conversion{goalId:%v,revenue:%v,negative:%v}",
		c.goalId, c.revenue, c.negative)
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
	if c.metadata != nil && len(c.metadata) > 0 {
		qb.Append(utils.QPMetadata, c.encodeMetadata())
	}
	return qb.String()
}

func (c *Conversion) encodeMetadata() string {
	sb := strings.Builder{}
	sb.WriteRune('{')
	addComma := false
	addedIndices := make(map[int]struct{})
	for _, mcd := range c.metadata {
		if mcd == nil {
			continue
		}
		if _, contains := addedIndices[mcd.ID()]; contains {
			continue
		}
		if addComma {
			sb.WriteRune(',')
		} else {
			addComma = true
		}
		writeCustomDataMetadata(mcd, &sb)
		addedIndices[mcd.ID()] = struct{}{}
	}
	sb.WriteRune('}')
	return sb.String()
}

func writeCustomDataMetadata(cd *CustomData, sb *strings.Builder) {
	sb.WriteRune('"')
	sb.WriteString(strconv.Itoa(cd.ID()))
	sb.WriteString("\":[")
	for i, value := range cd.Values() {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('"')
		value = strings.ReplaceAll(value, "\\", "\\\\")
		value = strings.ReplaceAll(value, "\"", "\\\"")
		sb.WriteString(value)
		sb.WriteRune('"')
	}
	sb.WriteRune(']')
}

func (c *Conversion) DataType() DataType {
	return DataTypeConversion
}
