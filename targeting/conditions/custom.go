package conditions

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/segmentio/encoding/json"

	"github.com/Kameleoon/client-go/types"
)

func NewCustomDatum(c types.TargetingCondition) *CustomDatum {
	include := false
	if c.Include != nil {
		include = *c.Include
	}
	return &CustomDatum{
		Type:     c.Type,
		Operator: c.Operator,
		Value:    c.Value,
		ID:       c.ID,
		Index:    c.Index,
		Weight:   c.Weight,
		Include:  include,
	}
}

type CustomDatum struct {
	Value    interface{}         `json:"value"`
	Type     types.TargetingType `json:"targetingType"`
	Operator types.OperatorType  `json:"valueMatchType"`
	ID       int                 `json:"id"`
	Index    string              `json:"customDataIndex"`
	Weight   int                 `json:"weight"`
	Include  bool                `json:"include"`
}

func (c *CustomDatum) String() string {
	if c == nil {
		return ""
	}
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	var s strings.Builder
	s.Grow(len(b))
	s.Write(b)
	return s.String()
}

func (c CustomDatum) GetType() types.TargetingType {
	return c.Type
}

func (c *CustomDatum) SetType(t types.TargetingType) {
	c.Type = t
}

func (c CustomDatum) GetInclude() bool {
	return c.Include
}

func (c *CustomDatum) SetInclude(i bool) {
	c.Include = i
}

func (c *CustomDatum) CheckTargeting(targetData []types.TargetingData) bool {
	var customData []*types.CustomData
	for _, td := range targetData {
		if td.Data.DataType() != types.DataTypeCustom {
			continue
		}
		custom, ok := td.Data.(*types.CustomData)
		if ok && custom.ID == c.Index {
			customData = append(customData, custom)
		}
	}
	if len(customData) == 0 {
		return c.Operator == types.OperatorUndefined
	}
	customDatum := customData[len(customData)-1]
	switch c.Operator {
	case types.OperatorContains:
		str, ok1 := customDatum.Value.(string)
		value, ok2 := c.Value.(string)
		if !ok1 || !ok2 {
			return false
		}
		return strings.Contains(str, value)
	case types.OperatorExact:
		if c.Value == customDatum.Value {
			return true
		}
	case types.OperatorMatch:
		str, ok1 := customDatum.Value.(string)
		pattern, ok2 := c.Value.(string)
		if !ok1 || !ok2 {
			return false
		}
		matched, err := regexp.MatchString(pattern, str)
		if err == nil && matched {
			return true
		}
	case types.OperatorLower, types.OperatorGreater, types.OperatorEqual:
		var number int
		switch v := c.Value.(type) {
		case string:
			number, _ = strconv.Atoi(v)
		case int:
			number = v
		default:
			return false
		}
		var value int
		switch v := customDatum.Value.(type) {
		case string:
			value, _ = strconv.Atoi(v)
		case int:
			value = v
		default:
			return false
		}
		switch c.Operator {
		case types.OperatorLower:
			if value < number {
				return true
			}
		case types.OperatorEqual:
			if value == number {
				return true
			}
		case types.OperatorGreater:
			if value > number {
				return true
			}
		}
	case types.OperatorIsTrue:
		if val, ok := customDatum.Value.(bool); ok && val {
			return true
		}
	case types.OperatorIsFalse:
		if val, ok := customDatum.Value.(bool); ok && !val {
			return true
		}
	}
	return false
}
